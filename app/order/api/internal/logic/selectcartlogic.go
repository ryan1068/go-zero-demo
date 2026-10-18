package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectCartLogic {
	return &SelectCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SelectCartLogic) SelectCart(req *types.SelectCartReq) (*types.BaseResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.SelectCart(ctx, &order.SelectCartReq{Selected: req.Selected})
	if err != nil {
		l.Errorf("select cart rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
