package logic

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"go-zero-demo/pkg/errcode"
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
	admin := req.Admin == 1
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

	where := []string{"p.deleted_at IS NULL"}
	args := make([]any, 0, 6)
	if !admin {
		where = append(where, "p.status = 1")
	} else {
		switch req.Status {
		case 1:
			where = append(where, "p.status = 1")
		case 2:
			where = append(where, "p.status = 0")
		}
	}
	if req.CategoryId > 0 {
		ids := []int64{req.CategoryId}
		childRows, err := l.svcCtx.DB.QueryContext(l.ctx,
			"SELECT id FROM category WHERE parent_id = ? AND deleted_at IS NULL", req.CategoryId,
		)
		if err != nil {
			l.Errorf("query child categories failed: %v", err)
			return nil, err
		}
		for childRows.Next() {
			var childID int64
			if err := childRows.Scan(&childID); err != nil {
				if closeErr := childRows.Close(); closeErr != nil {
					l.Errorf("close child category rows failed: %v", closeErr)
				}
				l.Errorf("scan child category failed: %v", err)
				return nil, err
			}
			ids = append(ids, childID)
		}
		if err := childRows.Err(); err != nil {
			if closeErr := childRows.Close(); closeErr != nil {
				l.Errorf("close child category rows failed: %v", closeErr)
			}
			l.Errorf("iterate child categories failed: %v", err)
			return nil, err
		}
		if err := childRows.Close(); err != nil {
			l.Errorf("close child category rows failed: %v", err)
			return nil, err
		}
		placeholders := make([]string, 0, len(ids))
		for _, id := range ids {
			placeholders = append(placeholders, "?")
			args = append(args, id)
		}
		where = append(where, "p.category_id IN ("+strings.Join(placeholders, ",")+")")
	}
	if keyword != "" {
		if utf8.RuneCountInString(keyword) >= 2 {
			where = append(where, "MATCH(p.name, p.subtitle, p.detail) AGAINST (? IN NATURAL LANGUAGE MODE)")
		} else {
			where = append(where, "p.name LIKE ?")
			keyword = "%" + keyword + "%"
		}
		args = append(args, keyword)
	}
	if req.MinPrice > 0 {
		where = append(where, "p.price >= ?")
		args = append(args, req.MinPrice)
	}
	if req.MaxPrice > 0 {
		where = append(where, "p.price <= ?")
		args = append(args, req.MaxPrice)
	}
	if req.LowStock == 1 {
		where = append(where, "p.stock <= ?")
		args = append(args, lowStockValue)
	}
	condition := strings.Join(where, " AND ")

	var total int64
	countQuery := "SELECT COUNT(*) FROM (SELECT 1 FROM product p WHERE " + condition +
		" LIMIT " + strconv.Itoa(maxProductCount) + ") AS filtered"
	if err := l.svcCtx.DB.QueryRowContext(l.ctx, countQuery, args...).Scan(&total); err != nil {
		l.Errorf("count products failed: %v", err)
		return nil, err
	}

	orderBy := "p.sort DESC, p.id DESC"
	switch req.Sort {
	case "sales":
		orderBy = "p.sales DESC, p.id DESC"
	case "price_asc":
		orderBy = "p.price ASC, p.id DESC"
	case "price_desc":
		orderBy = "p.price DESC, p.id DESC"
	case "newest":
		orderBy = "p.id DESC"
	}

	offset := int64(page-1) * int64(pageSize)
	listArgs := append(append([]any{}, args...), int64(pageSize), offset)
	query := "SELECT " + productSelectColumns + " FROM product p LEFT JOIN category c ON c.id = p.category_id AND c.deleted_at IS NULL WHERE " +
		condition + " ORDER BY " + orderBy + " LIMIT ? OFFSET ?"
	rows, err := l.svcCtx.DB.QueryContext(l.ctx, query, listArgs...)
	if err != nil {
		l.Errorf("query products failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close product rows failed: %v", closeErr)
		}
	}()

	list := make([]*product.ProductInfo, 0, pageSize)
	for rows.Next() {
		item, err := scanProduct(rows)
		if err != nil {
			l.Errorf("scan product failed: %v", err)
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate products failed: %v", err)
		return nil, err
	}

	return &product.ListProductsResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProductLogic) GetProduct(in *product.ProductReq) (*product.ProductResp, error) {
	id := in.Id
	admin := in.Admin == 1
	if id <= 0 {
		return &product.ProductResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	query := "SELECT " + productSelectColumns + " FROM product p LEFT JOIN category c ON c.id = p.category_id AND c.deleted_at IS NULL WHERE p.id = ? AND p.deleted_at IS NULL"
	if !admin {
		query += " AND p.status = 1"
	}
	info, err := scanProduct(l.svcCtx.DB.QueryRowContext(l.ctx, query, id))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &product.ProductResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query product failed: %v", err)
		return nil, err
	}
	return &product.ProductResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

type SaveProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveProductLogic {
	return &SaveProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveProductLogic) SaveProduct(req *product.SaveProductReq) (*product.ProductResp, error) {
	name := strings.TrimSpace(req.Name)
	subtitle := strings.TrimSpace(req.Subtitle)
	cover := strings.TrimSpace(req.CoverImage)
	unit := strings.TrimSpace(req.Unit)
	if unit == "" {
		unit = "件"
	}

	imagesList := make([]string, 0, len(req.Images))
	for _, img := range req.Images {
		img = strings.TrimSpace(img)
		if img == "" {
			continue
		}
		if !validText(img, maxCoverImage) {
			return &product.ProductResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
		}
		imagesList = append(imagesList, img)
	}
	if len(imagesList) > maxProductImages {
		return &product.ProductResp{Code: errcode.BadRequest, Msg: "商品图片不能超过 8 张"}, nil
	}
	if cover == "" && len(imagesList) > 0 {
		cover = imagesList[0]
	}

	if req.CategoryId <= 0 || name == "" || !validText(name, maxProductName) ||
		!validText(subtitle, maxSubtitle) || !validText(cover, maxCoverImage) ||
		!validText(unit, maxUnit) || !validText(req.Detail, maxDetailLength) ||
		req.Price <= 0 || req.OriginalPrice < 0 || req.Stock < 0 ||
		(req.Status != 0 && req.Status != 1) || (req.IsRecommend != 0 && req.IsRecommend != 1) {
		return &product.ProductResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var categoryCount int64
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM category WHERE id = ? AND deleted_at IS NULL", req.CategoryId,
	).Scan(&categoryCount); err != nil {
		l.Errorf("query category failed: %v", err)
		return nil, err
	}
	if categoryCount == 0 {
		return &product.ProductResp{Code: errcode.BadRequest, Msg: "分类不存在"}, nil
	}

	images, err := encodeImages(imagesList)
	if err != nil {
		l.Errorf("encode product images failed: %v", err)
		return nil, err
	}

	var id int64
	if req.Id > 0 {
		result, err := l.svcCtx.DB.ExecContext(l.ctx,
			"UPDATE product SET category_id = ?, name = ?, subtitle = ?, cover_image = ?, images = ?, detail = ?, price = ?, original_price = ?, stock = ?, unit = ?, status = ?, is_recommend = ?, sort = ? WHERE id = ? AND deleted_at IS NULL",
			req.CategoryId, name, subtitle, cover, images, req.Detail, req.Price, req.OriginalPrice, req.Stock, unit, req.Status, req.IsRecommend, req.Sort, req.Id,
		)
		if err != nil {
			l.Errorf("update product failed: %v", err)
			return nil, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			l.Errorf("read affected rows failed: %v", err)
			return nil, err
		}
		if affected == 0 {
			if _, err := l.findProduct(req.Id); errors.Is(err, sql.ErrNoRows) {
				return &product.ProductResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
			} else if err != nil {
				l.Errorf("query product failed: %v", err)
				return nil, err
			}
		}
		id = req.Id
	} else {
		result, err := l.svcCtx.DB.ExecContext(l.ctx,
			"INSERT INTO product (category_id, name, subtitle, cover_image, images, detail, price, original_price, stock, unit, status, is_recommend, sort) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			req.CategoryId, name, subtitle, cover, images, req.Detail, req.Price, req.OriginalPrice, req.Stock, unit, req.Status, req.IsRecommend, req.Sort,
		)
		if err != nil {
			l.Errorf("insert product failed: %v", err)
			return nil, err
		}
		id, err = result.LastInsertId()
		if err != nil {
			l.Errorf("read inserted product id failed: %v", err)
			return nil, err
		}
	}

	info, err := l.findProduct(id)
	if err != nil {
		l.Errorf("query saved product failed: %v", err)
		return nil, err
	}
	return &product.ProductResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

func (l *SaveProductLogic) findProduct(id int64) (*product.ProductInfo, error) {
	return scanProduct(l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT "+productSelectColumns+" FROM product p LEFT JOIN category c ON c.id = p.category_id AND c.deleted_at IS NULL WHERE p.id = ? AND p.deleted_at IS NULL", id,
	))
}

type DeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteProductLogic) DeleteProduct(in *product.ProductIdReq) (*product.BaseResp, error) {
	id := in.Id
	if id <= 0 {
		return &product.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	result, err := l.svcCtx.DB.ExecContext(l.ctx,
		"UPDATE product SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id,
	)
	if err != nil {
		l.Errorf("delete product failed: %v", err)
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		l.Errorf("read affected rows failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &product.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &product.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}

type UpdateProductStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductStatusLogic {
	return &UpdateProductStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateProductStatusLogic) UpdateProductStatus(req *product.UpdateProductStatusReq) (*product.BaseResp, error) {
	if req.Id <= 0 || (req.Status != 0 && req.Status != 1) {
		return &product.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	result, err := l.svcCtx.DB.ExecContext(l.ctx,
		"UPDATE product SET status = ? WHERE id = ? AND deleted_at IS NULL", req.Status, req.Id,
	)
	if err != nil {
		l.Errorf("update product status failed: %v", err)
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		l.Errorf("read affected rows failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &product.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &product.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}

type ProductStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductStatsLogic {
	return &ProductStatsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ProductStatsLogic) ProductStats(_ *product.EmptyReq) (*product.ProductStatsResp, error) {
	var stats product.ProductStats
	queries := []struct {
		query string
		dest  *int64
	}{
		{"SELECT COUNT(*) FROM product WHERE deleted_at IS NULL", &stats.Total},
		{"SELECT COUNT(*) FROM product WHERE deleted_at IS NULL AND status = 1", &stats.OnSale},
		{"SELECT COUNT(*) FROM product WHERE deleted_at IS NULL AND status = 0", &stats.OffSale},
		{"SELECT COUNT(*) FROM product WHERE deleted_at IS NULL AND status = 1 AND stock <= ?", &stats.LowStock},
	}
	for _, item := range queries {
		var err error
		if strings.Contains(item.query, "?") {
			err = l.svcCtx.DB.QueryRowContext(l.ctx, item.query, lowStockValue).Scan(item.dest)
		} else {
			err = l.svcCtx.DB.QueryRowContext(l.ctx, item.query).Scan(item.dest)
		}
		if err != nil {
			l.Errorf("query product stats failed: %v", err)
			return nil, err
		}
	}
	return &product.ProductStatsResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: &stats}, nil
}
