package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListOrdersLogic {
	return &AdminListOrdersLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminListOrdersLogic) AdminListOrders(req *types.ListOrdersReq) (*types.OrderListResp, error) {
	return NewListOrdersLogic(l.ctx, l.svcCtx).listOrders(req, true)
}
