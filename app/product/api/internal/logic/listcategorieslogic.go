package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
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

func (l *ListCategoriesLogic) ListCategories() (*types.ListCategoriesResp, error) {
	return l.listCategories(false)
}

func (l *ListCategoriesLogic) listCategories(admin bool) (*types.ListCategoriesResp, error) {
	s, err := l.svcCtx.ProductRpc.ListCategories(l.ctx, &product.ListCategoriesReq{Admin: adminFlag(admin)})
	if err != nil {
		l.Errorf("list categories rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ListCategoriesResp{Code: s.Code, Msg: s.Msg, Data: make([]types.CategoryInfo, 0, len(s.List))}
	for _, item := range s.List {
		if item != nil {
			resp.Data = append(resp.Data, toCategoryInfo(item))
		}
	}
	return resp, nil
}
