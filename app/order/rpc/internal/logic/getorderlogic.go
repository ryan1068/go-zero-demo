package logic

import (
	"context"
	"errors"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderReq) (*order.OrderResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	row, err := l.svcCtx.OrdersModel.FindOrder(l.ctx, nil, in.Id, userID, in.Admin == 1)
	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		return &order.OrderResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query order failed: %v", err)
		return nil, err
	}

	items, err := l.svcCtx.OrderItemModel.FindOrderItems(l.ctx, nil, in.Id)
	if err != nil {
		l.Errorf("query order items failed: %v", err)
		return nil, err
	}
	return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toOrderInfo(row, items)}, nil
}
