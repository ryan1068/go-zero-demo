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

func (l *ListCategoriesLogic) ListCategories(admin bool) (*types.ListCategoriesResp, error) {
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

type SaveCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveCategoryLogic {
	return &SaveCategoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveCategoryLogic) SaveCategory(req *types.SaveCategoryReq) (*types.CategoryResp, error) {
	s, err := l.svcCtx.ProductRpc.SaveCategory(l.ctx, &product.SaveCategoryReq{
		Id:       req.Id,
		ParentId: req.ParentId,
		Name:     req.Name,
		Icon:     req.Icon,
		Sort:     req.Sort,
		Status:   req.Status,
	})
	if err != nil {
		l.Errorf("save category rpc failed: %v", err)
		return nil, err
	}
	resp := &types.CategoryResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toCategoryInfo(s.Data)
	}
	return resp, nil
}

type DeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteCategoryLogic) DeleteCategory(id int64) (*types.BaseResp, error) {
	s, err := l.svcCtx.ProductRpc.DeleteCategory(l.ctx, &product.CategoryReq{Id: id})
	if err != nil {
		l.Errorf("delete category rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
