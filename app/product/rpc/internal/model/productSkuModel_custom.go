package model

import "context"

// ListByProduct returns active SKUs of a product.
func (m *customProductSkuModel) ListByProduct(ctx context.Context, productID int64) ([]*ProductSku, error) {
	var list []*ProductSku
	err := m.conn.QueryRowsCtx(ctx, &list,
		"select "+productSkuRows+" from "+m.table+
			" where `product_id` = ? and `deleted_at` is null order by `id` asc",
		productID,
	)
	return list, err
}

// FindActive returns a non-deleted SKU by id.
func (m *customProductSkuModel) FindActive(ctx context.Context, id uint64) (*ProductSku, error) {
	var resp ProductSku
	err := m.conn.QueryRowCtx(ctx, &resp,
		"select "+productSkuRows+" from "+m.table+" where `id` = ? and `deleted_at` is null limit 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// InsertSku inserts a SKU and returns its id.
func (m *customProductSkuModel) InsertSku(ctx context.Context, data *ProductSku) (uint64, error) {
	result, err := m.Insert(ctx, data)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

// UpdateSku updates the editable fields and returns affected rows.
func (m *customProductSkuModel) UpdateSku(ctx context.Context, data *ProductSku) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `product_id` = ?, `sku_code` = ?, `specs` = ?, `price` = ?,"+
			" `original_price` = ?, `stock` = ?, `status` = ?"+
			" where `id` = ? and `deleted_at` is null",
		data.ProductId, data.SkuCode, data.Specs, data.Price, data.OriginalPrice, data.Stock,
		data.Status, data.Id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// SoftDeleteSku marks a SKU as deleted.
func (m *customProductSkuModel) SoftDeleteSku(ctx context.Context, id uint64) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `deleted_at` = now() where `id` = ? and `deleted_at` is null", id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
