package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetAfterSaleLogic {
	return &AdminGetAfterSaleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminGetAfterSaleLogic) AdminGetAfterSale(req *types.AfterSaleReq) (*types.AfterSaleResp, error) {
	return NewGetAfterSaleLogic(l.ctx, l.svcCtx).get(req.Id, true)
}
