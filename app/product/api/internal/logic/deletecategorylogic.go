package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
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

func (l *DeleteCategoryLogic) DeleteCategory(req *types.PathIdReq) (*types.BaseResp, error) {
	s, err := l.svcCtx.ProductRpc.DeleteCategory(l.ctx, &product.CategoryReq{Id: req.Id})
	if err != nil {
		l.Errorf("delete category rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
