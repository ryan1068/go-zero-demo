package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCartLogic {
	return &DeleteCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteCartLogic) DeleteCart(req *types.CartItemReq) (*types.BaseResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.DeleteCart(ctx, &order.CartItemReq{Id: req.Id})
	if err != nil {
		l.Errorf("delete cart rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
