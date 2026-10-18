package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (*types.OrderResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.OrderResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.CreateOrder(ctx, &order.CreateOrderReq{
		AddressId: req.AddressId,
		Remark:    req.Remark,
		BuyNow:    req.BuyNow,
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	})
	if err != nil {
		l.Errorf("create order rpc failed: %v", err)
		return nil, err
	}
	resp := &types.OrderResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toOrderInfo(s.Data)
	}
	return resp, nil
}
