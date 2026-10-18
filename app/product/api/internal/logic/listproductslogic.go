package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductsLogic {
	return &ListProductsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListProductsLogic) ListProducts(req *types.ListProductsReq) (*types.ListProductsResp, error) {
	return l.listProducts(req, false)
}

func (l *ListProductsLogic) listProducts(req *types.ListProductsReq, admin bool) (*types.ListProductsResp, error) {
	s, err := l.svcCtx.ProductRpc.ListProducts(l.ctx, &product.ListProductsReq{
		Page:       req.Page,
		PageSize:   req.PageSize,
		CategoryId: req.CategoryId,
		Keyword:    req.Keyword,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Status:     req.Status,
		Sort:       req.Sort,
		LowStock:   req.LowStock,
		Admin:      adminFlag(admin),
	})
	if err != nil {
		l.Errorf("list products rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ListProductsResp{Code: s.Code, Msg: s.Msg}
	resp.Data.Total = s.Total
	resp.Data.Page = s.Page
	resp.Data.PageSize = s.PageSize
	resp.Data.List = make([]types.ProductInfo, 0, len(s.List))
	for _, item := range s.List {
		if item != nil {
			resp.Data.List = append(resp.Data.List, toProductInfo(item))
		}
	}
	return resp, nil
}
