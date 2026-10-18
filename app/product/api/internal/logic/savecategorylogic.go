package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

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
