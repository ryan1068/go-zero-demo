package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CrossProductModel accesses the product database from order-rpc. The stock
// deduction must stay in the same MySQL transaction as the order writes.
type CrossProductModel struct {
	conn sqlx.SqlConn
}

// NewCrossProductModel returns a cross-schema product model.
func NewCrossProductModel(conn sqlx.SqlConn) *CrossProductModel {
	return &CrossProductModel{conn: conn}
}

// ProductSnapshot is the product subset needed when creating an order.
type ProductSnapshot struct {
	Id     int64  `db:"id"`
	Name   string `db:"name"`
	Cover  string `db:"cover"`
	Price  int64  `db:"price"`
	Stock  int64  `db:"stock"`
	Status int64  `db:"status"`
}

// FindForUpdate locks a product row.
func (m *CrossProductModel) FindForUpdate(ctx context.Context, session sqlx.Session, productID int64) (*ProductSnapshot, error) {
	var row ProductSnapshot
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &row,
		"select `id`, `name`, coalesce(`cover_image`, '') as cover, `price`, `stock`, `status`"+
			" from looklook_product.product where `id` = ? and `deleted_at` is null for update",
		productID,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// DeductStock deducts stock and increases sales.
func (m *CrossProductModel) DeductStock(ctx context.Context, session sqlx.Session, productID, quantity int64) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update looklook_product.product set `stock` = `stock` - ?, `sales` = `sales` + ?"+
			" where `id` = ? and `stock` >= ?",
		quantity, quantity, productID, quantity,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// RestoreStock returns stock and decreases sales without going below zero.
func (m *CrossProductModel) RestoreStock(ctx context.Context, session sqlx.Session, productID, quantity int64) error {
	_, err := session.ExecCtx(ctx,
		"update looklook_product.product set `stock` = `stock` + ?, `sales` = greatest(`sales` - ?, 0)"+
			" where `id` = ?",
		quantity, quantity, productID,
	)
	return err
}

// CountProducts counts active products.
func (m *CrossProductModel) CountProducts(ctx context.Context, session sqlx.Session) (int64, error) {
	var total int64
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total,
		"select count(*) from looklook_product.product where `deleted_at` is null",
	)
	return total, err
}

// CountLowStock counts on-sale products at or below the threshold.
func (m *CrossProductModel) CountLowStock(ctx context.Context, session sqlx.Session, threshold int64) (int64, error) {
	var total int64
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total,
		"select count(*) from looklook_product.product"+
			" where `deleted_at` is null and `status` = 1 and `stock` <= ?", threshold,
	)
	return total, err
}
