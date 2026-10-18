package logic

import (
	"context"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListCategoriesLogic) ListCategories(in *product.ListCategoriesReq) (*product.ListCategoriesResp, error) {
	rows, err := l.svcCtx.CategoryModel.List(l.ctx, in.Admin == 1)
	if err != nil {
		l.Errorf("query categories failed: %v", err)
		return nil, err
	}

	list := make([]*product.CategoryInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCategoryInfo(row))
	}
	return &product.ListCategoriesResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), List: list}, nil
}
