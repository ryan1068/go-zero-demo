package logic

import (
	"encoding/json"
	"unicode/utf8"

	"go-zero-demo/product-rpc/internal/model"
	"go-zero-demo/product-rpc/pb/product"
)

const (
	defaultPageSize  = 12
	maxPageSize      = 100
	maxProductCount  = 10000
	lowStockValue    = 10
	maxProductImages = 8

	maxCategoryName = 64
	maxIconLength   = 512
	maxProductName  = 128
	maxSubtitle     = 255
	maxCoverImage   = 512
	maxUnit         = 16
	maxDetailLength = 20000
	maxSkuCode      = 64
	maxSkuSpecs     = 500
)

func validText(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func encodeImages(images []string) (string, error) {
	if images == nil {
		images = []string{}
	}
	data, err := json.Marshal(images)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func decodeImages(raw string) []string {
	images := make([]string, 0, 4)
	if raw == "" {
		return images
	}
	if err := json.Unmarshal([]byte(raw), &images); err != nil {
		return []string{}
	}
	return images
}

func toCategoryInfo(m *model.Category) *product.CategoryInfo {
	if m == nil {
		return nil
	}
	return &product.CategoryInfo{
		Id:        int64(m.Id),
		ParentId:  m.ParentId,
		Name:      m.Name,
		Icon:      m.Icon,
		Sort:      int32(m.Sort),
		Status:    int32(m.Status),
		CreatedAt: m.CreatedAt.Unix(),
		UpdatedAt: m.UpdatedAt.Unix(),
	}
}

func toProductInfo(row *model.ProductRow) *product.ProductInfo {
	if row == nil {
		return nil
	}
	return &product.ProductInfo{
		Id:            int64(row.Id),
		CategoryId:    row.CategoryId,
		CategoryName:  row.CategoryName,
		Name:          row.Name,
		Subtitle:      row.Subtitle,
		CoverImage:    row.CoverImage,
		Images:        decodeImages(row.Images),
		Detail:        row.Detail,
		Price:         row.Price,
		OriginalPrice: row.OriginalPrice,
		Stock:         int32(row.Stock),
		Sales:         int32(row.Sales),
		Unit:          row.Unit,
		Status:        int32(row.Status),
		IsRecommend:   int32(row.IsRecommend),
		Sort:          int32(row.Sort),
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func toSkuInfo(m *model.ProductSku) *product.SkuInfo {
	if m == nil {
		return nil
	}
	return &product.SkuInfo{
		Id:            int64(m.Id),
		ProductId:     m.ProductId,
		SkuCode:       m.SkuCode,
		Specs:         m.Specs,
		Price:         m.Price,
		OriginalPrice: m.OriginalPrice,
		Stock:         int32(m.Stock),
		Sales:         int32(m.Sales),
		Status:        int32(m.Status),
		CreatedAt:     m.CreatedAt.Unix(),
		UpdatedAt:     m.UpdatedAt.Unix(),
	}
}
