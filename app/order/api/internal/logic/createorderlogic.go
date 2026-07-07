package logic

import (
	"context"
	"go-zero-demo/order-rpc/pb/order"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CommonResp, err error) {
	s, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order.CreateOrderReq{
		GoodsId:   req.Goods_id,
		GoodsName: req.Goods_name,
		Price:     req.Price,
		CreatedAt: req.Created_at,
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
