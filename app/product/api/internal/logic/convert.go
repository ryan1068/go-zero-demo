package logic

import (
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"
)

func toCategoryInfo(info *product.CategoryInfo) types.CategoryInfo {
	return types.CategoryInfo{
		Id:        info.Id,
		ParentId:  info.ParentId,
		Name:      info.Name,
		Icon:      info.Icon,
		Sort:      info.Sort,
		Status:    info.Status,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}

func toProductInfo(info *product.ProductInfo) types.ProductInfo {
	return types.ProductInfo{
		Id:            info.Id,
		CategoryId:    info.CategoryId,
		CategoryName:  info.CategoryName,
		Name:          info.Name,
		Subtitle:      info.Subtitle,
		CoverImage:    info.CoverImage,
		Images:        info.Images,
		Detail:        info.Detail,
		Price:         info.Price,
		OriginalPrice: info.OriginalPrice,
		Stock:         info.Stock,
		Sales:         info.Sales,
		Unit:          info.Unit,
		Status:        info.Status,
		IsRecommend:   info.IsRecommend,
		Sort:          info.Sort,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}
}

func adminFlag(admin bool) int32 {
	if admin {
		return 1
	}
	return 0
}

func toSkuInfo(info *product.SkuInfo) types.SkuInfo {
	return types.SkuInfo{
		Id:            info.Id,
		ProductId:     info.ProductId,
		SkuCode:       info.SkuCode,
		Specs:         info.Specs,
		Price:         info.Price,
		OriginalPrice: info.OriginalPrice,
		Stock:         info.Stock,
		Sales:         info.Sales,
		Status:        info.Status,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}
}
