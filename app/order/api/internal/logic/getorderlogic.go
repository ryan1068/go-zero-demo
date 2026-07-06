package logic

import (
	"context"
	"go-zero-demo/order-rpc/pb/order"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderReq) (resp *types.CommonResp, err error) {
	s, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	resp = new(types.CommonResp)
	resp.Msg = s.Msg
	resp.Code = s.Code
	resp.Data = s.Data
	return
}
