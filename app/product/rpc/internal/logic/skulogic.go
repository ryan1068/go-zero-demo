package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	maxSkuCode  = 64
	maxSkuSpecs = 500
)

const skuSelectColumns = "id, product_id, sku_code, specs, price, original_price, stock, sales, status, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at)"

func scanSku(row rowScanner) (*product.SkuInfo, error) {
	info := new(product.SkuInfo)
	if err := row.Scan(
		&info.Id,
		&info.ProductId,
		&info.SkuCode,
		&info.Specs,
		&info.Price,
		&info.OriginalPrice,
		&info.Stock,
		&info.Sales,
		&info.Status,
		&info.CreatedAt,
		&info.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return info, nil
}

func loadSku(ctx context.Context, db *sql.DB, id int64) (*product.SkuInfo, error) {
	return scanSku(db.QueryRowContext(ctx,
		"SELECT "+skuSelectColumns+" FROM product_sku WHERE id = ? AND deleted_at IS NULL", id,
	))
}

type ListSkusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSkusLogic {
	return &ListSkusLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListSkusLogic) ListSkus(in *product.ProductIdReq) (*product.SkuListResp, error) {
	if in.Id <= 0 {
		return &product.SkuListResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	rows, err := l.svcCtx.DB.QueryContext(l.ctx,
		"SELECT "+skuSelectColumns+" FROM product_sku WHERE product_id = ? AND deleted_at IS NULL ORDER BY id ASC", in.Id,
	)
	if err != nil {
		l.Errorf("query skus failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close sku rows failed: %v", closeErr)
		}
	}()
	list := make([]*product.SkuInfo, 0, 8)
	for rows.Next() {
		item, err := scanSku(rows)
		if err != nil {
			l.Errorf("scan sku failed: %v", err)
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate skus failed: %v", err)
		return nil, err
	}
	return &product.SkuListResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), List: list}, nil
}

type SaveSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveSkuLogic {
	return &SaveSkuLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SaveSkuLogic) SaveSku(in *product.SaveSkuReq) (*product.SkuResp, error) {
	skuCode := strings.TrimSpace(in.SkuCode)
	specs := strings.TrimSpace(in.Specs)
	if specs == "" {
		specs = "{}"
	}
	if in.ProductId <= 0 || skuCode == "" || !validText(skuCode, maxSkuCode) ||
		!validText(specs, maxSkuSpecs) || in.Price <= 0 || in.OriginalPrice < 0 ||
		in.Stock < 0 || (in.Status != 0 && in.Status != 1) || !json.Valid([]byte(specs)) {
		return &product.SkuResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var productCount int64
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM product WHERE id = ? AND deleted_at IS NULL", in.ProductId,
	).Scan(&productCount); err != nil {
		l.Errorf("query product for sku failed: %v", err)
		return nil, err
	}
	if productCount == 0 {
		return &product.SkuResp{Code: errcode.BadRequest, Msg: "商品不存在"}, nil
	}

	var id int64
	if in.Id > 0 {
		result, err := l.svcCtx.DB.ExecContext(l.ctx,
			"UPDATE product_sku SET product_id = ?, sku_code = ?, specs = ?, price = ?, original_price = ?, stock = ?, status = ? WHERE id = ? AND deleted_at IS NULL",
			in.ProductId, skuCode, specs, in.Price, in.OriginalPrice, in.Stock, in.Status, in.Id,
		)
		if err != nil {
			l.Errorf("update sku failed: %v", err)
			return nil, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			l.Errorf("read sku affected rows failed: %v", err)
			return nil, err
		}
		if affected == 0 {
			if _, err := loadSku(l.ctx, l.svcCtx.DB, in.Id); errors.Is(err, sql.ErrNoRows) {
				return &product.SkuResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
			} else if err != nil {
				l.Errorf("query sku failed: %v", err)
				return nil, err
			}
		}
		id = in.Id
	} else {
		result, err := l.svcCtx.DB.ExecContext(l.ctx,
			"INSERT INTO product_sku (product_id, sku_code, specs, price, original_price, stock, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
			in.ProductId, skuCode, specs, in.Price, in.OriginalPrice, in.Stock, in.Status,
		)
		if err != nil {
			l.Errorf("insert sku failed: %v", err)
			return nil, err
		}
		id, err = result.LastInsertId()
		if err != nil {
			l.Errorf("read sku id failed: %v", err)
			return nil, err
		}
	}

	info, err := loadSku(l.ctx, l.svcCtx.DB, id)
	if err != nil {
		l.Errorf("query saved sku failed: %v", err)
		return nil, err
	}
	return &product.SkuResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

type DeleteSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSkuLogic {
	return &DeleteSkuLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteSkuLogic) DeleteSku(in *product.SkuIdReq) (*product.BaseResp, error) {
	if in.Id <= 0 {
		return &product.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	result, err := l.svcCtx.DB.ExecContext(l.ctx,
		"UPDATE product_sku SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", in.Id,
	)
	if err != nil {
		l.Errorf("delete sku failed: %v", err)
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		l.Errorf("read sku affected rows failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &product.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &product.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
