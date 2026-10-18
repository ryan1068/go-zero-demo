package logic

import (
	"context"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteCategoryLogic) DeleteCategory(in *product.CategoryReq) (*product.BaseResp, error) {
	id := in.Id
	if id <= 0 {
		return &product.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	count, err := l.svcCtx.CategoryModel.CountProducts(l.ctx, uint64(id))
	if err != nil {
		l.Errorf("count category products failed: %v", err)
		return nil, err
	}
	if count > 0 {
		return &product.BaseResp{Code: errcode.Conflict, Msg: "分类下还有商品，不能删除"}, nil
	}

	childCount, err := l.svcCtx.CategoryModel.CountChildren(l.ctx, uint64(id))
	if err != nil {
		l.Errorf("count child categories failed: %v", err)
		return nil, err
	}
	if childCount > 0 {
		return &product.BaseResp{Code: errcode.Conflict, Msg: "分类下还有子分类，不能删除"}, nil
	}

	affected, err := l.svcCtx.CategoryModel.SoftDeleteCategory(l.ctx, uint64(id))
	if err != nil {
		l.Errorf("delete category failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &product.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &product.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
