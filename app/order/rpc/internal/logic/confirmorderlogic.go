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

type ConfirmOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmOrderLogic {
	return &ConfirmOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConfirmOrderLogic) ConfirmOrder(in *order.ConfirmOrderReq) (*order.OrderResp, error) {
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

		if lock.Status == orderStatusCompleted {
			info, err := loadOrderInfo(ctx, l.svcCtx, session, in.Id, userID, false)
			if err != nil {
				return err
			}
			resp = &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}
			return nil
		}
		if lock.Status != orderStatusShipped {
			resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "订单未发货，不能确认收货"}
			return model.ErrAbortTransaction
		}

		if _, err := l.svcCtx.OrdersModel.UpdateOrderCompletedByUser(ctx, session, in.Id, userID, orderStatusCompleted); err != nil {
			return err
		}
		if err := l.svcCtx.OrderLogModel.InsertOrderLog(ctx, session, model.NewOrderLog{
			OrderID:      in.Id,
			FromStatus:   orderStatusShipped,
			ToStatus:     orderStatusCompleted,
			OperatorType: operatorUser,
			OperatorID:   userID,
			Remark:       "确认收货",
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
		l.Errorf("confirm order failed: %v", err)
		return nil, err
	}
	return resp, nil
}
