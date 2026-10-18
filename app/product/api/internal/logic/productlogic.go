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

func (l *ListProductsLogic) ListProducts(req *types.ListProductsReq, admin bool) (*types.ListProductsResp, error) {
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

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProductLogic) GetProduct(id int64, admin bool) (*types.ProductResp, error) {
	s, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.ProductReq{Id: id, Admin: adminFlag(admin)})
	if err != nil {
		l.Errorf("get product rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ProductResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toProductInfo(s.Data)
	}
	return resp, nil
}

type SaveProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveProductLogic {
	return &SaveProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveProductLogic) SaveProduct(req *types.SaveProductReq) (*types.ProductResp, error) {
	s, err := l.svcCtx.ProductRpc.SaveProduct(l.ctx, &product.SaveProductReq{
		Id:            req.Id,
		CategoryId:    req.CategoryId,
		Name:          req.Name,
		Subtitle:      req.Subtitle,
		CoverImage:    req.CoverImage,
		Images:        req.Images,
		Detail:        req.Detail,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		Unit:          req.Unit,
		Status:        req.Status,
		IsRecommend:   req.IsRecommend,
		Sort:          req.Sort,
	})
	if err != nil {
		l.Errorf("save product rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ProductResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toProductInfo(s.Data)
	}
	return resp, nil
}

type DeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteProductLogic) DeleteProduct(id int64) (*types.BaseResp, error) {
	s, err := l.svcCtx.ProductRpc.DeleteProduct(l.ctx, &product.ProductIdReq{Id: id})
	if err != nil {
		l.Errorf("delete product rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}

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

type ProductStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductStatsLogic {
	return &ProductStatsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ProductStatsLogic) ProductStats() (*types.ProductStatsResp, error) {
	s, err := l.svcCtx.ProductRpc.ProductStats(l.ctx, &product.EmptyReq{})
	if err != nil {
		l.Errorf("product stats rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ProductStatsResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = types.ProductStats{
			Total:    s.Data.Total,
			OnSale:   s.Data.OnSale,
			OffSale:  s.Data.OffSale,
			LowStock: s.Data.LowStock,
		}
	}
	return resp, nil
}
