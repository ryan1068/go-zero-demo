package logic

import (
	"context"
	"errors"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *PayOrderLogic) PayOrder(in *order.PayOrderReq) (*order.OrderResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var resp *order.OrderResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		lock, err := l.svcCtx.OrdersModel.LockOrder(ctx, session, in.Id, &userID)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.OrderResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}

		switch lock.Status {
		case orderStatusPendingShipment, orderStatusShipped, orderStatusCompleted:
			info, err := loadOrderInfo(ctx, l.svcCtx, session, in.Id, userID, false)
			if err != nil {
				return err
			}
			resp = &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}
			return nil
		case orderStatusCanceled:
			resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "订单已取消"}
			return model.ErrAbortTransaction
		case orderStatusPendingPayment:
			// Continue below.
		default:
			resp = &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}
			return model.ErrAbortTransaction
		}

		transactionID := "MOCK-" + generateNo("")
		affected, err := l.svcCtx.OrdersModel.UpdateOrderPaid(ctx, session, in.Id, userID,
			orderStatusPendingPayment, orderStatusPendingShipment, payTypeMock, transactionID)
		if err != nil {
			return err
		}
		if affected == 0 {
			resp = &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}
			return model.ErrAbortTransaction
		}
		if err := l.svcCtx.PaymentRecordModel.InsertPaymentRecord(ctx, session, model.NewPaymentRecord{
			PaymentNo:     generateNo("P"),
			OrderID:       in.Id,
			OrderNo:       lock.OrderNo,
			UserID:        userID,
			Amount:        lock.PayAmount,
			PayType:       payTypeMock,
			Status:        1,
			TransactionID: transactionID,
		}); err != nil {
			return err
		}
		if err := l.svcCtx.OrderLogModel.InsertOrderLog(ctx, session, model.NewOrderLog{
			OrderID:      in.Id,
			FromStatus:   orderStatusPendingPayment,
			ToStatus:     orderStatusPendingShipment,
			OperatorType: operatorUser,
			OperatorID:   userID,
			Remark:       "模拟支付成功",
		}); err != nil {
			return err
		}

		info, err := loadOrderInfo(ctx, l.svcCtx, session, in.Id, userID, false)
		if err != nil {
			return err
		}
		resp = &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}
		return nil
	})
	switch {
	case errors.Is(err, model.ErrAbortTransaction):
		return resp, nil
	case err != nil:
		l.Errorf("pay order failed: %v", err)
		return nil, err
	}
	return resp, nil
}
