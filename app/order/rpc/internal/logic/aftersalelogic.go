package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"unicode/utf8"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	afterSaleStatusPending        = 10
	afterSaleStatusApproved       = 20
	afterSaleStatusRejected       = 30
	afterSaleStatusBuyerShipped   = 40
	afterSaleStatusSellerReceived = 50
	afterSaleStatusRefunded       = 60
	afterSaleStatusClosed         = 70

	maxAfterSaleReason      = 255
	maxAfterSaleDescription = 500
	maxAfterSaleImages      = 8
	maxAfterSaleImageLength = 512
	maxAfterSaleRemark      = 255
	maxAfterSaleTrackingNo  = 64
)

const afterSaleSelectColumns = `a.id, a.after_sale_no, a.order_id, a.order_item_id, a.user_id, a.type,
	a.reason, a.description, COALESCE(CAST(a.images AS CHAR), '[]'), a.amount, a.status,
	a.admin_remark, a.tracking_no, UNIX_TIMESTAMP(a.apply_at), COALESCE(UNIX_TIMESTAMP(a.review_at), 0),
	COALESCE(UNIX_TIMESTAMP(a.buyer_shipped_at), 0), COALESCE(UNIX_TIMESTAMP(a.seller_received_at), 0),
	COALESCE(UNIX_TIMESTAMP(a.refunded_at), 0), UNIX_TIMESTAMP(a.created_at), UNIX_TIMESTAMP(a.updated_at),
	COALESCE(o.order_no, ''), COALESCE(i.product_name, ''), COALESCE(i.product_cover, '')`

type afterSaleRowScanner interface {
	Scan(dest ...any) error
}

func scanAfterSale(row afterSaleRowScanner) (*order.AfterSaleInfo, error) {
	info := new(order.AfterSaleInfo)
	var imagesJSON string
	if err := row.Scan(
		&info.Id,
		&info.AfterSaleNo,
		&info.OrderId,
		&info.OrderItemId,
		&info.UserId,
		&info.Type,
		&info.Reason,
		&info.Description,
		&imagesJSON,
		&info.Amount,
		&info.Status,
		&info.AdminRemark,
		&info.TrackingNo,
		&info.ApplyAt,
		&info.ReviewAt,
		&info.BuyerShippedAt,
		&info.SellerReceivedAt,
		&info.RefundedAt,
		&info.CreatedAt,
		&info.UpdatedAt,
		&info.OrderNo,
		&info.ProductName,
		&info.ProductCover,
	); err != nil {
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

func loadAfterSale(ctx context.Context, db *sql.DB, id, userID int64, admin bool) (*order.AfterSaleInfo, error) {
	query := "SELECT " + afterSaleSelectColumns +
		" FROM after_sale a LEFT JOIN orders o ON o.id = a.order_id LEFT JOIN order_item i ON i.id = a.order_item_id" +
		" WHERE a.id = ? AND a.deleted_at IS NULL"
	args := []any{id}
	if !admin {
		query += " AND a.user_id = ?"
		args = append(args, userID)
	}
	return scanAfterSale(db.QueryRowContext(ctx, query, args...))
}

type ApplyAfterSaleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAfterSaleLogic {
	return &ApplyAfterSaleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ApplyAfterSaleLogic) ApplyAfterSale(in *order.ApplyAfterSaleReq) (*order.AfterSaleResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.AfterSaleResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	reason := strings.TrimSpace(in.Reason)
	description := strings.TrimSpace(in.Description)
	if in.OrderId <= 0 || (in.Type != 1 && in.Type != 2) || reason == "" ||
		!validAfterSaleText(reason, maxAfterSaleReason) || !validAfterSaleText(description, maxAfterSaleDescription) ||
		len(in.Images) > maxAfterSaleImages {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	images := make([]string, 0, len(in.Images))
	for _, image := range in.Images {
		image = strings.TrimSpace(image)
		if image == "" {
			continue
		}
		if !validAfterSaleText(image, maxAfterSaleImageLength) {
			return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
		}
		images = append(images, image)
	}
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		l.Errorf("encode after sale images failed: %v", err)
		return nil, err
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin after sale transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	var orderUserID int64
	var orderStatus int32
	err = tx.QueryRowContext(l.ctx,
		"SELECT user_id, status FROM orders WHERE id = ? AND deleted_at IS NULL FOR UPDATE", in.OrderId,
	).Scan(&orderUserID, &orderStatus)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query order for after sale failed: %v", err)
		return nil, err
	}
	if orderUserID != userID {
		return &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	if orderStatus != orderStatusPendingShipment && orderStatus != orderStatusShipped && orderStatus != orderStatusCompleted {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前订单状态不支持售后"}, nil
	}

	var itemID, price int64
	var quantity int32
	itemQuery := "SELECT id, price, quantity FROM order_item WHERE order_id = ?"
	itemArgs := []any{in.OrderId}
	if in.OrderItemId > 0 {
		itemQuery += " AND id = ?"
		itemArgs = append(itemArgs, in.OrderItemId)
	}
	itemQuery += " ORDER BY id LIMIT 1"
	err = tx.QueryRowContext(l.ctx, itemQuery, itemArgs...).Scan(&itemID, &price, &quantity)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.AfterSaleResp{Code: errcode.NotFound, Msg: "订单明细不存在"}, nil
	case err != nil:
		l.Errorf("query order item for after sale failed: %v", err)
		return nil, err
	}

	var active int64
	if err := tx.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM after_sale WHERE order_item_id = ? AND user_id = ? AND status NOT IN (?, ?, ?) AND deleted_at IS NULL",
		itemID, userID, afterSaleStatusRejected, afterSaleStatusRefunded, afterSaleStatusClosed,
	).Scan(&active); err != nil {
		l.Errorf("query active after sale failed: %v", err)
		return nil, err
	}
	if active > 0 {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "该商品已有进行中的售后"}, nil
	}

	afterSaleNo := generateNo("AS")
	amount := price * int64(quantity)
	result, err := tx.ExecContext(l.ctx,
		"INSERT INTO after_sale (after_sale_no, order_id, order_item_id, user_id, type, reason, description, images, amount, status, apply_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())",
		afterSaleNo, in.OrderId, itemID, userID, in.Type, reason, description, string(imagesJSON), amount, afterSaleStatusPending,
	)
	if err != nil {
		l.Errorf("insert after sale failed: %v", err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		l.Errorf("read after sale id failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit after sale transaction failed: %v", err)
		return nil, err
	}

	info, err := loadAfterSale(l.ctx, l.svcCtx.DB, id, userID, false)
	if err != nil {
		l.Errorf("query saved after sale failed: %v", err)
		return nil, err
	}
	return &order.AfterSaleResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

type ListAfterSalesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAfterSalesLogic {
	return &ListAfterSalesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListAfterSalesLogic) ListAfterSales(in *order.ListAfterSalesReq) (*order.AfterSaleListResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.AfterSaleListResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	where := []string{"a.deleted_at IS NULL"}
	args := make([]any, 0, 3)
	if in.Admin != 1 {
		where = append(where, "a.user_id = ?")
		args = append(args, userID)
	}
	if in.Status > 0 {
		where = append(where, "a.status = ?")
		args = append(args, in.Status)
	}
	condition := strings.Join(where, " AND ")

	var total int64
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM after_sale a WHERE "+condition, args...,
	).Scan(&total); err != nil {
		l.Errorf("count after sales failed: %v", err)
		return nil, err
	}

	offset := int64(page-1) * int64(pageSize)
	listArgs := append(append([]any{}, args...), int64(pageSize), offset)
	rows, err := l.svcCtx.DB.QueryContext(l.ctx,
		"SELECT "+afterSaleSelectColumns+
			" FROM after_sale a LEFT JOIN orders o ON o.id = a.order_id LEFT JOIN order_item i ON i.id = a.order_item_id"+
			" WHERE "+condition+" ORDER BY a.id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	if err != nil {
		l.Errorf("query after sales failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close after sale rows failed: %v", closeErr)
		}
	}()

	list := make([]*order.AfterSaleInfo, 0, pageSize)
	for rows.Next() {
		item, err := scanAfterSale(rows)
		if err != nil {
			l.Errorf("scan after sale failed: %v", err)
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate after sales failed: %v", err)
		return nil, err
	}
	return &order.AfterSaleListResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

type GetAfterSaleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAfterSaleLogic {
	return &GetAfterSaleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetAfterSaleLogic) GetAfterSale(in *order.AfterSaleReq) (*order.AfterSaleResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.AfterSaleResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	info, err := loadAfterSale(l.ctx, l.svcCtx.DB, in.Id, userID, in.Admin == 1)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query after sale failed: %v", err)
		return nil, err
	}
	return &order.AfterSaleResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

type HandleAfterSaleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleAfterSaleLogic {
	return &HandleAfterSaleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *HandleAfterSaleLogic) HandleAfterSale(in *order.HandleAfterSaleReq) (*order.AfterSaleResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.AfterSaleResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	action := strings.ToLower(strings.TrimSpace(in.Action))
	adminRemark := strings.TrimSpace(in.AdminRemark)
	trackingNo := strings.TrimSpace(in.TrackingNo)
	if in.Id <= 0 || !validAfterSaleText(adminRemark, maxAfterSaleRemark) || !validAfterSaleText(trackingNo, maxAfterSaleTrackingNo) {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin after sale handle transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	var ownerID, amount, orderID int64
	var saleType, status int32
	err = tx.QueryRowContext(l.ctx,
		"SELECT user_id, type, status, amount, order_id FROM after_sale WHERE id = ? AND deleted_at IS NULL FOR UPDATE", in.Id,
	).Scan(&ownerID, &saleType, &status, &amount, &orderID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query after sale for handle failed: %v", err)
		return nil, err
	}
	if in.Admin != 1 && ownerID != userID {
		return &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	if in.Admin != 1 && action != "ship" {
		return &order.AfterSaleResp{Code: errcode.Unauthorized, Msg: "需要管理员权限"}, nil
	}
	if in.Admin == 1 && action == "ship" {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "寄回操作由买家发起"}, nil
	}

	switch action {
	case "approve":
		if status != afterSaleStatusPending {
			return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可审核"}, nil
		}
		if _, err := tx.ExecContext(l.ctx,
			"UPDATE after_sale SET status = ?, review_at = NOW(), admin_remark = ? WHERE id = ?",
			afterSaleStatusApproved, adminRemark, in.Id,
		); err != nil {
			l.Errorf("approve after sale failed: %v", err)
			return nil, err
		}
	case "reject":
		if status != afterSaleStatusPending {
			return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可审核"}, nil
		}
		if _, err := tx.ExecContext(l.ctx,
			"UPDATE after_sale SET status = ?, review_at = NOW(), admin_remark = ? WHERE id = ?",
			afterSaleStatusRejected, adminRemark, in.Id,
		); err != nil {
			l.Errorf("reject after sale failed: %v", err)
			return nil, err
		}
	case "ship":
		if status != afterSaleStatusApproved || saleType != 2 {
			return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可寄回"}, nil
		}
		if trackingNo == "" {
			return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "请填写退回物流单号"}, nil
		}
		if _, err := tx.ExecContext(l.ctx,
			"UPDATE after_sale SET status = ?, tracking_no = ?, buyer_shipped_at = NOW() WHERE id = ?",
			afterSaleStatusBuyerShipped, trackingNo, in.Id,
		); err != nil {
			l.Errorf("ship after sale failed: %v", err)
			return nil, err
		}
	case "receive":
		if status != afterSaleStatusBuyerShipped {
			return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可收货"}, nil
		}
		if _, err := tx.ExecContext(l.ctx,
			"UPDATE after_sale SET status = ?, seller_received_at = NOW(), admin_remark = ? WHERE id = ?",
			afterSaleStatusSellerReceived, adminRemark, in.Id,
		); err != nil {
			l.Errorf("receive after sale failed: %v", err)
			return nil, err
		}
	case "refund":
		if status != afterSaleStatusApproved && status != afterSaleStatusSellerReceived {
			return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可退款"}, nil
		}
		if _, err := tx.ExecContext(l.ctx,
			"UPDATE after_sale SET status = ?, refunded_at = NOW(), admin_remark = ? WHERE id = ?",
			afterSaleStatusRefunded, adminRemark, in.Id,
		); err != nil {
			l.Errorf("refund after sale failed: %v", err)
			return nil, err
		}
		if _, err := tx.ExecContext(l.ctx,
			"INSERT INTO refund_record (refund_no, after_sale_id, order_id, user_id, amount, channel, status, refunded_at) VALUES (?, ?, ?, ?, ?, 'mock', 1, NOW())",
			generateNo("R"), in.Id, orderID, ownerID, amount,
		); err != nil {
			l.Errorf("insert refund record failed: %v", err)
			return nil, err
		}
	default:
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "不支持的售后操作"}, nil
	}

	if err := tx.Commit(); err != nil {
		l.Errorf("commit after sale handle transaction failed: %v", err)
		return nil, err
	}
	info, err := loadAfterSale(l.ctx, l.svcCtx.DB, in.Id, userID, true)
	if err != nil {
		l.Errorf("query handled after sale failed: %v", err)
		return nil, err
	}
	return &order.AfterSaleResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

func validAfterSaleText(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}
