package logic

import (
	"encoding/json"
	"unicode/utf8"

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

	categorySelectColumns = "id, parent_id, name, COALESCE(icon, ''), sort, status, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at)"
	productSelectColumns  = "p.id, p.category_id, COALESCE(c.name, ''), p.name, COALESCE(p.subtitle, ''), COALESCE(p.cover_image, ''), COALESCE(CAST(p.images AS CHAR), '[]'), COALESCE(p.detail, ''), p.price, p.original_price, p.stock, p.sales, COALESCE(p.unit, ''), p.status, p.is_recommend, p.sort, UNIX_TIMESTAMP(p.created_at), UNIX_TIMESTAMP(p.updated_at)"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCategory(row rowScanner) (*product.CategoryInfo, error) {
	info := new(product.CategoryInfo)
	err := row.Scan(
		&info.Id,
		&info.ParentId,
		&info.Name,
		&info.Icon,
		&info.Sort,
		&info.Status,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func scanProduct(row rowScanner) (*product.ProductInfo, error) {
	info := new(product.ProductInfo)
	var imagesJSON string
	err := row.Scan(
		&info.Id,
		&info.CategoryId,
		&info.CategoryName,
		&info.Name,
		&info.Subtitle,
		&info.CoverImage,
		&imagesJSON,
		&info.Detail,
		&info.Price,
		&info.OriginalPrice,
		&info.Stock,
		&info.Sales,
		&info.Unit,
		&info.Status,
		&info.IsRecommend,
		&info.Sort,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	info.Images = make([]string, 0, 4)
	if imagesJSON != "" {
		if err := json.Unmarshal([]byte(imagesJSON), &info.Images); err != nil {
			info.Images = []string{}
		}
	}
	return info, nil
}

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
