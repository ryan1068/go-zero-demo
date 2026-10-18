package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAfterSaleLogic {
	return &ApplyAfterSaleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ApplyAfterSaleLogic) ApplyAfterSale(req *types.ApplyAfterSaleReq) (*types.AfterSaleResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AfterSaleResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.ApplyAfterSale(ctx, &order.ApplyAfterSaleReq{
		OrderId:     req.OrderId,
		OrderItemId: req.OrderItemId,
		Type:        req.Type,
		Reason:      req.Reason,
		Description: req.Description,
		Images:      req.Images,
	})
	if err != nil {
		l.Errorf("apply after sale rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AfterSaleResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAfterSaleInfo(s.Data)
	}
	return resp, nil
}
