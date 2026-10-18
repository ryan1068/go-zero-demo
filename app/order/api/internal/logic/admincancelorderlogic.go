package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCancelOrderLogic {
	return &AdminCancelOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminCancelOrderLogic) AdminCancelOrder(req *types.CancelOrderReq) (*types.OrderResp, error) {
	return NewCancelOrderLogic(l.ctx, l.svcCtx).cancel(req, true)
}
