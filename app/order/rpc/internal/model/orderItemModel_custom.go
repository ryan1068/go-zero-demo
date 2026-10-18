package model

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// OrderItemRow is an order item snapshot.
type OrderItemRow struct {
	Id           int64  `db:"id"`
	ProductId    int64  `db:"product_id"`
	ProductName  string `db:"product_name"`
	ProductCover string `db:"product_cover"`
	Price        int64  `db:"price"`
	Quantity     int64  `db:"quantity"`
	TotalAmount  int64  `db:"total_amount"`
}

// AfterSaleItem is the order item subset used when applying after-sales.
type AfterSaleItem struct {
	Id       int64 `db:"id"`
	Price    int64 `db:"price"`
	Quantity int64 `db:"quantity"`
}

// StockLine is one product quantity in an order.
type StockLine struct {
	ProductId int64 `db:"product_id"`
	Quantity  int64 `db:"quantity"`
}

// NewOrderItem is one item to insert.
type NewOrderItem struct {
	ProductID    int64
	ProductName  string
	ProductCover string
	Price        int64
	Quantity     int64
	TotalAmount  int64
}

// FindOrderItems returns the items of an order.
func (m *customOrderItemModel) FindOrderItems(ctx context.Context, session sqlx.Session, orderID int64) ([]*OrderItemRow, error) {
	var list []*OrderItemRow
	err := sessionOrConn(m.conn, session).QueryRowsCtx(ctx, &list,
		"select `id`, `product_id`, `product_name`, `product_cover`, `price`, `quantity`, `total_amount`"+
			" from "+m.table+" where `order_id` = ? order by `id`", orderID,
	)
	return list, err
}

// FindItemForAfterSale returns the target item, or the first item when id is 0.
func (m *customOrderItemModel) FindItemForAfterSale(ctx context.Context, session sqlx.Session, orderID, itemID int64) (*AfterSaleItem, error) {
	query := "select `id`, `price`, `quantity` from " + m.table + " where `order_id` = ?"
	args := []any{orderID}
	if itemID > 0 {
		query += " and `id` = ?"
		args = append(args, itemID)
	}
	query += " order by `id` limit 1"
	var item AfterSaleItem
	if err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &item, query, args...); err != nil {
		return nil, err
	}
	return &item, nil
}

// InsertOrderItems inserts all items of an order.
func (m *customOrderItemModel) InsertOrderItems(ctx context.Context, session sqlx.Session, orderID int64, items []NewOrderItem) error {
	for _, item := range items {
		if _, err := session.ExecCtx(ctx,
			"insert into "+m.table+" (`order_id`, `product_id`, `product_name`, `product_cover`,"+
				" `price`, `quantity`, `total_amount`) values (?, ?, ?, ?, ?, ?, ?)",
			orderID, item.ProductID, item.ProductName, item.ProductCover, item.Price,
			item.Quantity, item.TotalAmount,
		); err != nil {
			return err
		}
	}
	return nil
}

// StockLines returns the product quantities of an order.
func (m *customOrderItemModel) StockLines(ctx context.Context, session sqlx.Session, orderID int64) ([]*StockLine, error) {
	var lines []*StockLine
	err := sessionOrConn(m.conn, session).QueryRowsCtx(ctx, &lines,
		"select `product_id`, `quantity` from "+m.table+" where `order_id` = ?", orderID,
	)
	return lines, err
}

// OrderItemWithOrder is an order item carrying its order id for batch loading.
type OrderItemWithOrder struct {
	OrderId      int64  `db:"order_id"`
	Id           int64  `db:"id"`
	ProductId    int64  `db:"product_id"`
	ProductName  string `db:"product_name"`
	ProductCover string `db:"product_cover"`
	Price        int64  `db:"price"`
	Quantity     int64  `db:"quantity"`
	TotalAmount  int64  `db:"total_amount"`
}

// FindOrderItemsByOrderIDs loads items for several orders in one query.
func (m *customOrderItemModel) FindOrderItemsByOrderIDs(ctx context.Context, session sqlx.Session, orderIDs []int64) (map[int64][]*OrderItemRow, error) {
	result := make(map[int64][]*OrderItemRow, len(orderIDs))
	if len(orderIDs) == 0 {
		return result, nil
	}
	placeholders := make([]string, 0, len(orderIDs))
	args := make([]any, 0, len(orderIDs))
	for _, id := range orderIDs {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}
	var rows []*OrderItemWithOrder
	query := "select `order_id`, `id`, `product_id`, `product_name`, `product_cover`, `price`," +
		" `quantity`, `total_amount` from " + m.table + " where `order_id` in (" +
		strings.Join(placeholders, ",") + ") order by `order_id`, `id`"
	if err := sessionOrConn(m.conn, session).QueryRowsCtx(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.OrderId] = append(result[row.OrderId], &OrderItemRow{
			Id:           row.Id,
			ProductId:    row.ProductId,
			ProductName:  row.ProductName,
			ProductCover: row.ProductCover,
			Price:        row.Price,
			Quantity:     row.Quantity,
			TotalAmount:  row.TotalAmount,
		})
	}
	return result, nil
}
