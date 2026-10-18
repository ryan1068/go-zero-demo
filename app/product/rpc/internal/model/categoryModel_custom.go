package model

import "context"

// List returns categories ordered for display.
func (m *customCategoryModel) List(ctx context.Context, admin bool) ([]*Category, error) {
	query := "select " + categoryRows + " from " + m.table + " where `deleted_at` is null"
	if !admin {
		query += " and `status` = 1"
	}
	query += " order by `sort` desc, `id` asc"
	var list []*Category
	err := m.conn.QueryRowsCtx(ctx, &list, query)
	return list, err
}

// FindActive returns a non-deleted category by id.
func (m *customCategoryModel) FindActive(ctx context.Context, id uint64) (*Category, error) {
	var resp Category
	err := m.conn.QueryRowCtx(ctx, &resp,
		"select "+categoryRows+" from "+m.table+" where `id` = ? and `deleted_at` is null limit 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// InsertCategory inserts a category and returns its id.
func (m *customCategoryModel) InsertCategory(ctx context.Context, data *Category) (uint64, error) {
	result, err := m.Insert(ctx, data)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

// UpdateCategory updates the editable fields and returns affected rows.
func (m *customCategoryModel) UpdateCategory(ctx context.Context, data *Category) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `parent_id` = ?, `name` = ?, `icon` = ?, `sort` = ?, `status` = ?"+
			" where `id` = ? and `deleted_at` is null",
		data.ParentId, data.Name, data.Icon, data.Sort, data.Status, data.Id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// SoftDeleteCategory marks a category as deleted.
func (m *customCategoryModel) SoftDeleteCategory(ctx context.Context, id uint64) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `deleted_at` = now() where `id` = ? and `deleted_at` is null", id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// CountProducts counts active products in a category.
func (m *customCategoryModel) CountProducts(ctx context.Context, id uint64) (int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total,
		"select count(*) from `product` where `category_id` = ? and `deleted_at` is null", id,
	)
	return total, err
}

// CountChildren counts active child categories.
func (m *customCategoryModel) CountChildren(ctx context.Context, id uint64) (int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total,
		"select count(*) from "+m.table+" where `parent_id` = ? and `deleted_at` is null", id,
	)
	return total, err
}

// Exists reports whether an active category exists.
func (m *customCategoryModel) Exists(ctx context.Context, id uint64) (bool, error) {
	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total,
		"select count(*) from "+m.table+" where `id` = ? and `deleted_at` is null", id,
	); err != nil {
		return false, err
	}
	return total > 0, nil
}
