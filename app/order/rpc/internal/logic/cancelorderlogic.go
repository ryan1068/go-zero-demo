package logic

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderReq) (*order.OrderResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	reason := strings.TrimSpace(in.Reason)
	if utf8.RuneCountInString(reason) > maxRemarkLength {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	admin := in.Admin == 1
	if reason == "" {
		if admin {
			reason = "管理员取消"
		} else {
			reason = "用户取消"
		}
	}

	var resp *order.OrderResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var owner *int64
		if !admin {
			owner = &userID
		}
		lock, err := l.svcCtx.OrdersModel.LockOrder(ctx, session, in.Id, owner)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.OrderResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}

		if lock.Status == orderStatusCanceled {
			info, err := loadOrderInfo(ctx, l.svcCtx, session, in.Id, userID, admin)
			if err != nil {
				return err
			}
			resp = &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}
			return nil
		}
		if lock.Status != orderStatusPendingPayment && lock.Status != orderStatusPendingShipment {
			resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "订单已发货，不能取消"}
			return model.ErrAbortTransaction
		}

		if _, err := l.svcCtx.OrdersModel.UpdateOrderCanceled(ctx, session, in.Id, orderStatusCanceled, reason); err != nil {
			return err
		}
		lines, err := l.svcCtx.OrderItemModel.StockLines(ctx, session, in.Id)
		if err != nil {
			return err
		}
		for _, line := range lines {
			if err := l.svcCtx.ProductModel.RestoreStock(ctx, session, line.ProductId, line.Quantity); err != nil {
				return err
			}
		}
		operatorType := int64(operatorUser)
		if admin {
			operatorType = operatorAdmin
		}
		if err := l.svcCtx.OrderLogModel.InsertOrderLog(ctx, session, model.NewOrderLog{
			OrderID:      in.Id,
			FromStatus:   lock.Status,
			ToStatus:     orderStatusCanceled,
			OperatorType: operatorType,
			OperatorID:   userID,
			Remark:       reason,
		}); err != nil {
			return err
		}

		info, err := loadOrderInfo(ctx, l.svcCtx, session, in.Id, userID, admin)
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
		l.Errorf("cancel order failed: %v", err)
		return nil, err
	}
	return resp, nil
}
