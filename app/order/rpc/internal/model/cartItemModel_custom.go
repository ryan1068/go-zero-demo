package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const cartSelectColumns = "ci.`id`, ci.`product_id`, coalesce(p.`name`, '') as product_name," +
	" coalesce(p.`cover_image`, '') as product_cover, coalesce(p.`price`, 0) as price," +
	" coalesce(p.`stock`, 0) as stock, ci.`quantity`, ci.`selected`," +
	" coalesce(p.`status`, 0) as product_status," +
	" unix_timestamp(ci.`created_at`) as created_at, unix_timestamp(ci.`updated_at`) as updated_at"

const cartFromClause = " from `cart_item` ci left join looklook_product.product p" +
	" on p.`id` = ci.`product_id` and p.`deleted_at` is null"

// CartItemRow is a cart item joined with its product.
type CartItemRow struct {
	Id            int64  `db:"id"`
	ProductId     int64  `db:"product_id"`
	ProductName   string `db:"product_name"`
	ProductCover  string `db:"product_cover"`
	Price         int64  `db:"price"`
	Stock         int64  `db:"stock"`
	Quantity      int64  `db:"quantity"`
	Selected      int64  `db:"selected"`
	ProductStatus int64  `db:"product_status"`
	CreatedAt     int64  `db:"created_at"`
	UpdatedAt     int64  `db:"updated_at"`
}

// CartSelection is a selected cart item used at checkout.
type CartSelection struct {
	ProductId int64 `db:"product_id"`
	Quantity  int64 `db:"quantity"`
}

// CartItemLock is the minimal cart row used under FOR UPDATE.
type CartItemLock struct {
	Id       int64 `db:"id"`
	Quantity int64 `db:"quantity"`
}

// CartProductLock is the product stock/status seen from a cart row.
type CartProductLock struct {
	Stock  int64 `db:"stock"`
	Status int64 `db:"status"`
}

// NewCartItem is a cart item to insert.
type NewCartItem struct {
	UserID    int64
	ProductID int64
	Quantity  int64
}

// FindSelectedForUpdate locks the selected cart items of a user.
func (m *customCartItemModel) FindSelectedForUpdate(ctx context.Context, session sqlx.Session, userID int64) ([]*CartSelection, error) {
	var list []*CartSelection
	err := session.QueryRowsCtx(ctx, &list,
		"select `product_id`, `quantity` from "+m.table+
			" where `user_id` = ? and `selected` = 1 for update", userID,
	)
	return list, err
}

// FindForUpdate locks one cart item by user and product.
func (m *customCartItemModel) FindForUpdate(ctx context.Context, session sqlx.Session, userID, productID int64) (*CartItemLock, error) {
	var item CartItemLock
	err := session.QueryRowCtx(ctx, &item,
		"select `id`, `quantity` from "+m.table+
			" where `user_id` = ? and `product_id` = ? for update", userID, productID,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// LockCartProduct locks a cart item and returns its product stock/status.
func (m *customCartItemModel) LockCartProduct(ctx context.Context, session sqlx.Session, id, userID int64) (*CartProductLock, error) {
	var lock CartProductLock
	err := session.QueryRowCtx(ctx, &lock,
		"select coalesce(p.`stock`, 0) as stock, coalesce(p.`status`, 0) as status"+
			" from `cart_item` ci left join looklook_product.product p"+
			" on p.`id` = ci.`product_id` and p.`deleted_at` is null"+
			" where ci.`id` = ? and ci.`user_id` = ? for update", id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &lock, nil
}

// InsertCartItem inserts a cart item and returns its id.
func (m *customCartItemModel) InsertCartItem(ctx context.Context, session sqlx.Session, data NewCartItem) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"insert into "+m.table+" (`user_id`, `product_id`, `quantity`, `selected`) values (?, ?, ?, 1)",
		data.UserID, data.ProductID, data.Quantity,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateCartItem updates quantity/selection of one cart item.
func (m *customCartItemModel) UpdateCartItem(ctx context.Context, session sqlx.Session, id, userID, quantity, selected int64) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `quantity` = ?, `selected` = ? where `id` = ? and `user_id` = ?",
		quantity, selected, id, userID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// DeleteCartItem deletes one cart item.
func (m *customCartItemModel) DeleteCartItem(ctx context.Context, session sqlx.Session, id, userID int64) (int64, error) {
	result, err := sessionOrConn(m.conn, session).ExecCtx(ctx,
		"delete from "+m.table+" where `id` = ? and `user_id` = ?", id, userID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateCartSelectedByUser updates the selection of all cart items of a user.
func (m *customCartItemModel) UpdateCartSelectedByUser(ctx context.Context, session sqlx.Session, userID, selected int64) error {
	_, err := sessionOrConn(m.conn, session).ExecCtx(ctx,
		"update "+m.table+" set `selected` = ? where `user_id` = ?", selected, userID,
	)
	return err
}

// DeleteSelectedCartItems removes checked-out cart items.
func (m *customCartItemModel) DeleteSelectedCartItems(ctx context.Context, session sqlx.Session, userID int64) error {
	_, err := session.ExecCtx(ctx,
		"delete from "+m.table+" where `user_id` = ? and `selected` = 1", userID,
	)
	return err
}

// FindCartItemRow returns one cart item joined with its product.
func (m *customCartItemModel) FindCartItemRow(ctx context.Context, session sqlx.Session, id, userID int64) (*CartItemRow, error) {
	var row CartItemRow
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &row,
		"select "+cartSelectColumns+cartFromClause+" where ci.`id` = ? and ci.`user_id` = ?",
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// ListCartItemRows returns the cart of a user.
func (m *customCartItemModel) ListCartItemRows(ctx context.Context, session sqlx.Session, userID int64) ([]*CartItemRow, error) {
	var list []*CartItemRow
	err := sessionOrConn(m.conn, session).QueryRowsCtx(ctx, &list,
		"select "+cartSelectColumns+cartFromClause+
			" where ci.`user_id` = ? order by ci.`updated_at` desc", userID,
	)
	return list, err
}
