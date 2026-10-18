package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListAfterSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListAfterSalesLogic {
	return &AdminListAfterSalesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminListAfterSalesLogic) AdminListAfterSales(req *types.ListAfterSalesReq) (*types.AfterSaleListResp, error) {
	return NewListAfterSalesLogic(l.ctx, l.svcCtx).list(req, true)
}
