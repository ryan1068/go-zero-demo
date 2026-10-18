package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductStatusLogic {
	return &UpdateProductStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateProductStatusLogic) UpdateProductStatus(req *types.UpdateProductStatusReq) (*types.BaseResp, error) {
	s, err := l.svcCtx.ProductRpc.UpdateProductStatus(l.ctx, &product.UpdateProductStatusReq{
		Id:     req.Id,
		Status: req.Status,
	})
	if err != nil {
		l.Errorf("update product status rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
