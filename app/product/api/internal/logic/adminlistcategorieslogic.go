package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListCategoriesLogic {
	return &AdminListCategoriesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminListCategoriesLogic) AdminListCategories() (*types.ListCategoriesResp, error) {
	return NewListCategoriesLogic(l.ctx, l.svcCtx).listCategories(true)
}
