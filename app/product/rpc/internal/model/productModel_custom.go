package model

import (
	"context"
	"strconv"
	"strings"
)

// ProductRow is a product joined with its category, using unix timestamps.
type ProductRow struct {
	Id            uint64 `db:"id"`
	CategoryId    int64  `db:"category_id"`
	CategoryName  string `db:"category_name"`
	Name          string `db:"name"`
	Subtitle      string `db:"subtitle"`
	CoverImage    string `db:"cover_image"`
	Images        string `db:"images"`
	Detail        string `db:"detail"`
	Price         int64  `db:"price"`
	OriginalPrice int64  `db:"original_price"`
	Stock         int64  `db:"stock"`
	Sales         int64  `db:"sales"`
	Unit          string `db:"unit"`
	Status        int64  `db:"status"`
	IsRecommend   int64  `db:"is_recommend"`
	Sort          int64  `db:"sort"`
	CreatedAt     int64  `db:"created_at"`
	UpdatedAt     int64  `db:"updated_at"`
}

// ProductFilter is the filter for the product list.
type ProductFilter struct {
	Admin      bool
	CategoryId int64
	Status     int64
	Keyword    string
	MinPrice   int64
	MaxPrice   int64
	LowStock   bool
	Sort       string
	Page       int64
	PageSize   int64
	MaxRows    int64
	LowStockAt int64
}

// ProductStats holds dashboard counters.
type ProductStats struct {
	Total    int64
	OnSale   int64
	OffSale  int64
	LowStock int64
}

const productSelect = "p.`id`, p.`category_id`, coalesce(c.`name`, '') as category_name," +
	" p.`name`, coalesce(p.`subtitle`, '') as subtitle, coalesce(p.`cover_image`, '') as cover_image," +
	" coalesce(cast(p.`images` as char), '[]') as images, coalesce(p.`detail`, '') as detail," +
	" p.`price`, p.`original_price`, p.`stock`, p.`sales`, coalesce(p.`unit`, '') as unit," +
	" p.`status`, p.`is_recommend`, p.`sort`," +
	" unix_timestamp(p.`created_at`) as created_at, unix_timestamp(p.`updated_at`) as updated_at"

const productFrom = " from `product` p left join `category` c" +
	" on c.`id` = p.`category_id` and c.`deleted_at` is null"

func productCondition(f ProductFilter) (string, []any) {
	where := []string{"p.`deleted_at` is null"}
	args := make([]any, 0, 6)
	if f.Admin {
		switch f.Status {
		case 1:
			where = append(where, "p.`status` = 1")
		case 2:
			where = append(where, "p.`status` = 0")
		}
	} else {
		where = append(where, "p.`status` = 1")
	}
	if f.CategoryId > 0 {
		where = append(where, "(p.`category_id` = ? or p.`category_id` in"+
			" (select `id` from `category` where `parent_id` = ? and `deleted_at` is null))")
		args = append(args, f.CategoryId, f.CategoryId)
	}
	if f.Keyword != "" {
		if len([]rune(f.Keyword)) >= 2 {
			where = append(where, "match(p.`name`, p.`subtitle`, p.`detail`) against (? in natural language mode)")
		} else {
			where = append(where, "p.`name` like ?")
			f.Keyword = "%" + f.Keyword + "%"
		}
		args = append(args, f.Keyword)
	}
	if f.MinPrice > 0 {
		where = append(where, "p.`price` >= ?")
		args = append(args, f.MinPrice)
	}
	if f.MaxPrice > 0 {
		where = append(where, "p.`price` <= ?")
		args = append(args, f.MaxPrice)
	}
	if f.LowStock {
		where = append(where, "p.`stock` <= ?")
		args = append(args, f.LowStockAt)
	}
	return strings.Join(where, " and "), args
}

func productOrderBy(sort string) string {
	switch sort {
	case "sales":
		return "p.`sales` desc, p.`id` desc"
	case "price_asc":
		return "p.`price` asc, p.`id` desc"
	case "price_desc":
		return "p.`price` desc, p.`id` desc"
	case "newest":
		return "p.`id` desc"
	default:
		return "p.`sort` desc, p.`id` desc"
	}
}

// ListProducts returns a page of products and the total count.
func (m *customProductModel) ListProducts(ctx context.Context, f ProductFilter) ([]*ProductRow, int64, error) {
	condition, args := productCondition(f)

	var total int64
	countQuery := "select count(*) from (select 1 from `product` p where " + condition +
		" limit " + strconv.FormatInt(f.MaxRows, 10) + ") as filtered"
	if err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	listArgs := append(append([]any{}, args...), f.PageSize, (f.Page-1)*f.PageSize)
	var list []*ProductRow
	query := "select " + productSelect + productFrom + " where " + condition +
		" order by " + productOrderBy(f.Sort) + " limit ? offset ?"
	if err := m.conn.QueryRowsCtx(ctx, &list, query, listArgs...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// FindProductRow returns one product joined with its category.
func (m *customProductModel) FindProductRow(ctx context.Context, id uint64, admin bool) (*ProductRow, error) {
	query := "select " + productSelect + productFrom + " where p.`id` = ? and p.`deleted_at` is null"
	if !admin {
		query += " and p.`status` = 1"
	}
	var row ProductRow
	if err := m.conn.QueryRowCtx(ctx, &row, query, id); err != nil {
		return nil, err
	}
	return &row, nil
}

// InsertProduct inserts a product and returns its id.
func (m *customProductModel) InsertProduct(ctx context.Context, data *Product) (uint64, error) {
	result, err := m.Insert(ctx, data)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

// UpdateProduct updates the editable fields and returns affected rows.
func (m *customProductModel) UpdateProduct(ctx context.Context, data *Product) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `category_id` = ?, `name` = ?, `subtitle` = ?, `cover_image` = ?,"+
			" `images` = ?, `detail` = ?, `price` = ?, `original_price` = ?, `stock` = ?,"+
			" `unit` = ?, `status` = ?, `is_recommend` = ?, `sort` = ?"+
			" where `id` = ? and `deleted_at` is null",
		data.CategoryId, data.Name, data.Subtitle, data.CoverImage, data.Images, data.Detail,
		data.Price, data.OriginalPrice, data.Stock, data.Unit, data.Status, data.IsRecommend,
		data.Sort, data.Id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// SoftDeleteProduct marks a product as deleted.
func (m *customProductModel) SoftDeleteProduct(ctx context.Context, id uint64) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `deleted_at` = now() where `id` = ? and `deleted_at` is null", id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateProductStatus updates the on-sale status.
func (m *customProductModel) UpdateProductStatus(ctx context.Context, id uint64, status int64) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `status` = ? where `id` = ? and `deleted_at` is null", status, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Exists reports whether an active product exists.
func (m *customProductModel) Exists(ctx context.Context, id uint64) (bool, error) {
	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total,
		"select count(*) from "+m.table+" where `id` = ? and `deleted_at` is null", id,
	); err != nil {
		return false, err
	}
	return total > 0, nil
}

// Stats returns product counters for the admin dashboard.
func (m *customProductModel) Stats(ctx context.Context, lowStock int64) (ProductStats, error) {
	var stats ProductStats
	queries := []struct {
		query string
		args  []any
		dest  *int64
	}{
		{"select count(*) from " + m.table + " where `deleted_at` is null", nil, &stats.Total},
		{"select count(*) from " + m.table + " where `deleted_at` is null and `status` = 1", nil, &stats.OnSale},
		{"select count(*) from " + m.table + " where `deleted_at` is null and `status` = 0", nil, &stats.OffSale},
		{"select count(*) from " + m.table + " where `deleted_at` is null and `status` = 1 and `stock` <= ?", []any{lowStock}, &stats.LowStock},
	}
	for _, item := range queries {
		if err := m.conn.QueryRowCtx(ctx, item.dest, item.query, item.args...); err != nil {
			return stats, err
		}
	}
	return stats, nil
}
