package logic

import (
	"context"
	"strings"
	"unicode/utf8"

	"go-zero-demo/order-rpc/internal/model"
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
	filter := model.OrderFilter{
		Admin:     admin,
		OrderNo:   orderNo,
		Status:    int64(in.Status),
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
		Page:      int64(page),
		PageSize:  int64(pageSize),
		MaxRows:   maxOrderCount,
	}
	if !admin {
		filter.UserID = userID
	} else if in.UserId > 0 {
		filter.UserID = in.UserId
	}

	total, err := l.svcCtx.OrdersModel.CountOrders(l.ctx, nil, filter)
	if err != nil {
		l.Errorf("count orders failed: %v", err)
		return nil, err
	}

	rows, err := l.svcCtx.OrdersModel.ListOrders(l.ctx, nil, filter)
	if err != nil {
		l.Errorf("query orders failed: %v", err)
		return nil, err
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.Id)
	}
	itemsByOrder, err := l.svcCtx.OrderItemModel.FindOrderItemsByOrderIDs(l.ctx, nil, ids)
	if err != nil {
		l.Errorf("query order items failed: %v", err)
		return nil, err
	}

	list := make([]*order.OrderInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, toOrderInfo(row, itemsByOrder[row.Id]))
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
