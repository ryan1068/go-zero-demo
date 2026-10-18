package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmOrderLogic {
	return &ConfirmOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ConfirmOrderLogic) ConfirmOrder(req *types.ConfirmOrderReq) (*types.OrderResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.OrderResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.ConfirmOrder(ctx, &order.ConfirmOrderReq{Id: req.Id})
	if err != nil {
		l.Errorf("confirm order rpc failed: %v", err)
		return nil, err
	}
	resp := &types.OrderResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toOrderInfo(s.Data)
	}
	return resp, nil
}
