package logic

import (
	"context"
	"strconv"
	"strings"
	"unicode/utf8"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListOrdersLogic) ListOrders(in *order.ListOrdersReq) (*order.OrderListResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderListResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}

	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	orderNo := strings.TrimSpace(in.OrderNo)
	if utf8.RuneCountInString(orderNo) > 32 {
		return &order.OrderListResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	if in.StartTime > 0 && in.EndTime > 0 && in.StartTime > in.EndTime {
		return &order.OrderListResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	admin := in.Admin == 1
	where := []string{"deleted_at IS NULL"}
	args := make([]any, 0, 6)
	if !admin {
		where = append(where, "user_id = ?")
		args = append(args, userID)
	} else if in.UserId > 0 {
		where = append(where, "user_id = ?")
		args = append(args, in.UserId)
	}
	if orderNo != "" {
		where = append(where, "order_no = ?")
		args = append(args, orderNo)
	}
	if in.Status > 0 {
		where = append(where, "status = ?")
		args = append(args, in.Status)
	}
	if in.StartTime > 0 {
		where = append(where, "created_at >= FROM_UNIXTIME(?)")
		args = append(args, in.StartTime)
	}
	if in.EndTime > 0 {
		where = append(where, "created_at <= FROM_UNIXTIME(?)")
		args = append(args, in.EndTime)
	}
	condition := strings.Join(where, " AND ")

	var total int64
	countQuery := "SELECT COUNT(*) FROM (SELECT 1 FROM orders WHERE " + condition +
		" LIMIT " + strconv.Itoa(maxOrderCount) + ") AS filtered"
	if err := l.svcCtx.DB.QueryRowContext(l.ctx, countQuery, args...).Scan(&total); err != nil {
		l.Errorf("count orders failed: %v", err)
		return nil, err
	}

	offset := int64(page-1) * int64(pageSize)
	listArgs := append(append([]any{}, args...), int64(pageSize), offset)
	rows, err := l.svcCtx.DB.QueryContext(l.ctx,
		"SELECT "+orderSelectColumns+" FROM orders WHERE "+condition+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	if err != nil {
		l.Errorf("query orders failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close order rows failed: %v", closeErr)
		}
	}()

	list := make([]*order.OrderInfo, 0, pageSize)
	for rows.Next() {
		info, err := scanOrderHeader(rows)
		if err != nil {
			l.Errorf("scan order failed: %v", err)
			return nil, err
		}
		list = append(list, info)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate orders failed: %v", err)
		return nil, err
	}
	for _, info := range list {
		items, err := loadOrderItems(l.ctx, l.svcCtx.DB, info.Id)
		if err != nil {
			l.Errorf("query order items failed: %v", err)
			return nil, err
		}
		info.Items = items
	}

	return &order.OrderListResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
