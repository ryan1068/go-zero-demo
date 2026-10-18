package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCartLogic {
	return &ListCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListCartLogic) ListCart() (*types.CartListResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.CartListResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.ListCart(ctx, &order.EmptyReq{})
	if err != nil {
		l.Errorf("list cart rpc failed: %v", err)
		return nil, err
	}
	resp := &types.CartListResp{Code: s.Code, Msg: s.Msg}
	resp.Data.Total = s.Total
	resp.Data.List = make([]types.CartItemInfo, 0, len(s.List))
	for _, item := range s.List {
		if item != nil {
			resp.Data.List = append(resp.Data.List, toCartItem(item))
		}
	}
	return resp, nil
}
