package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetOrderTracksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetOrderTracksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetOrderTracksLogic {
	return &AdminGetOrderTracksLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminGetOrderTracksLogic) AdminGetOrderTracks(req *types.GetOrderReq) (*types.OrderTracksResp, error) {
	return NewGetOrderTracksLogic(l.ctx, l.svcCtx).tracks(req.Id, true)
}
