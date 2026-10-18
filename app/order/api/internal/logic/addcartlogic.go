package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCartLogic {
	return &AddCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddCartLogic) AddCart(req *types.AddCartReq) (*types.CartItemResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.CartItemResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.AddCart(ctx, &order.AddCartReq{
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	})
	if err != nil {
		l.Errorf("add cart rpc failed: %v", err)
		return nil, err
	}
	resp := &types.CartItemResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toCartItem(s.Data)
	}
	return resp, nil
}
