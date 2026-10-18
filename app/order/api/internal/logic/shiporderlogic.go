package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipOrderLogic {
	return &ShipOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ShipOrderLogic) ShipOrder(req *types.ShipOrderReq) (*types.OrderResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.OrderResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.ShipOrder(ctx, &order.ShipOrderReq{
		Id:               req.Id,
		LogisticsCompany: req.LogisticsCompany,
		TrackingNo:       req.TrackingNo,
	})
	if err != nil {
		l.Errorf("ship order rpc failed: %v", err)
		return nil, err
	}
	resp := &types.OrderResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toOrderInfo(s.Data)
	}
	return resp, nil
}
