package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderReq) (*types.OrderResp, error) {
	return l.cancel(req, false)
}

func (l *CancelOrderLogic) cancel(req *types.CancelOrderReq, admin bool) (*types.OrderResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.OrderResp{Code: 401, Msg: "未授权"}, nil
	}
	adminFlag := int32(0)
	if admin {
		adminFlag = 1
	}
	s, err := l.svcCtx.OrderRpc.CancelOrder(ctx, &order.CancelOrderReq{
		Id:     req.Id,
		Reason: req.Reason,
		Admin:  adminFlag,
	})
	if err != nil {
		l.Errorf("cancel order rpc failed: %v", err)
		return nil, err
	}
	resp := &types.OrderResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toOrderInfo(s.Data)
	}
	return resp, nil
}
