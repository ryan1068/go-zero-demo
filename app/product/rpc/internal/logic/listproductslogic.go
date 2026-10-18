package logic

import (
	"context"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/model"
	"go-zero-demo/product-rpc/internal/svc"
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

func (l *ListProductsLogic) ListProducts(req *product.ListProductsReq) (*product.ListProductsResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	keyword := strings.TrimSpace(req.Keyword)
	if !validText(keyword, maxProductName) || req.MinPrice < 0 || req.MaxPrice < 0 ||
		(req.MinPrice > 0 && req.MaxPrice > 0 && req.MinPrice > req.MaxPrice) {
		return &product.ListProductsResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	list, total, err := l.svcCtx.ProductModel.ListProducts(l.ctx, model.ProductFilter{
		Admin:      req.Admin == 1,
		CategoryId: req.CategoryId,
		Status:     int64(req.Status),
		Keyword:    keyword,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		LowStock:   req.LowStock == 1,
		Sort:       req.Sort,
		Page:       int64(page),
		PageSize:   int64(pageSize),
		MaxRows:    maxProductCount,
		LowStockAt: lowStockValue,
	})
	if err != nil {
		l.Errorf("query products failed: %v", err)
		return nil, err
	}

	items := make([]*product.ProductInfo, 0, len(list))
	for _, row := range list {
		items = append(items, toProductInfo(row))
	}

	return &product.ListProductsResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
