package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetOrderLogic {
	return &AdminGetOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminGetOrderLogic) AdminGetOrder(req *types.GetOrderReq) (*types.OrderResp, error) {
	return NewGetOrderLogic(l.ctx, l.svcCtx).get(req.Id, true)
}
