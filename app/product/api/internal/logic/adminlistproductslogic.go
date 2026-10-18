package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListProductsLogic {
	return &AdminListProductsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminListProductsLogic) AdminListProducts(req *types.ListProductsReq) (*types.ListProductsResp, error) {
	return NewListProductsLogic(l.ctx, l.svcCtx).listProducts(req, true)
}
