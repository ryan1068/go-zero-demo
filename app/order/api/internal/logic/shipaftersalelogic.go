package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShipAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipAfterSaleLogic {
	return &ShipAfterSaleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ShipAfterSaleLogic) ShipAfterSale(req *types.ShipAfterSaleReq) (*types.AfterSaleResp, error) {
	return NewHandleAfterSaleLogic(l.ctx, l.svcCtx).handle(req.Id, "ship", "", req.TrackingNo, false)
}
