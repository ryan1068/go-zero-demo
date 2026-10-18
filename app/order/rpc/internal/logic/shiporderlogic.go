package logic

import (
	"context"
	"errors"
	"strings"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ShipOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipOrderLogic {
	return &ShipOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ShipOrderLogic) ShipOrder(in *order.ShipOrderReq) (*order.OrderResp, error) {
	adminID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	company := strings.TrimSpace(in.LogisticsCompany)
	trackingNo := strings.TrimSpace(in.TrackingNo)
	if in.Id <= 0 || company == "" || !trackNoValid(company, trackingNo) {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var resp *order.OrderResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		lock, err := l.svcCtx.OrdersModel.LockOrder(ctx, session, in.Id, nil)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.OrderResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}

		if lock.Status == orderStatusShipped {
			info, err := loadOrderInfo(ctx, l.svcCtx, session, in.Id, 0, true)
			if err != nil {
				return err
			}
			resp = &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}
			return nil
		}
		if lock.Status != orderStatusPendingShipment {
			resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "当前订单状态不能发货"}
			return model.ErrAbortTransaction
		}

		if _, err := l.svcCtx.OrdersModel.UpdateOrderShipped(ctx, session, in.Id, orderStatusShipped, company, trackingNo); err != nil {
			return err
		}
		if err := l.svcCtx.OrderLogModel.InsertOrderLog(ctx, session, model.NewOrderLog{
			OrderID:      in.Id,
			FromStatus:   orderStatusPendingShipment,
			ToStatus:     orderStatusShipped,
			OperatorType: operatorAdmin,
			OperatorID:   adminID,
			Remark:       "管理员发货",
		}); err != nil {
			return err
		}

		info, err := loadOrderInfo(ctx, l.svcCtx, session, in.Id, 0, true)
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
		l.Errorf("ship order failed: %v", err)
		return nil, err
	}
	return resp, nil
}
