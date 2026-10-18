package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	orderStatusPendingPayment  = 10
	orderStatusPendingShipment = 20
	orderStatusShipped         = 30
	orderStatusCompleted       = 40
	orderStatusCanceled        = 50

	payTypeMock = 1

	operatorUser   = 1
	operatorAdmin  = 2
	operatorSystem = 3

	defaultPageSize = 10
	maxPageSize     = 100
	maxOrderCount   = 10000

	maxQuantity         = 99
	maxRemarkLength     = 255
	maxReceiverName     = 64
	maxReceiverPhone    = 32
	maxReceiverRegion   = 64
	maxReceiverAddress  = 255
	maxReceiverPostal   = 16
	maxLogisticsCompany = 64
	maxTrackingNo       = 64
	lowStockThreshold   = 10

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

var noSequence uint64

func generateNo(prefix string) string {
	seq := atomic.AddUint64(&noSequence, 1)
	return fmt.Sprintf("%s%d%05d%03d", prefix, time.Now().UnixNano(), seq%100000, rand.IntN(1000))
}

func trackNoValid(company, trackingNo string) bool {
	return trackingNo != "" &&
		utf8.RuneCountInString(company) <= maxLogisticsCompany &&
		utf8.RuneCountInString(trackingNo) <= maxTrackingNo
}

func validAfterSaleText(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func decodeImages(raw string) []string {
	images := make([]string, 0, 4)
	if raw == "" {
		return images
	}
	if err := json.Unmarshal([]byte(raw), &images); err != nil {
		return []string{}
	}
	return images
}

func toOrderItemInfo(row *model.OrderItemRow) *order.OrderItemInfo {
	return &order.OrderItemInfo{
		Id:           row.Id,
		ProductId:    row.ProductId,
		ProductName:  row.ProductName,
		ProductCover: row.ProductCover,
		Price:        row.Price,
		Quantity:     int32(row.Quantity),
		TotalAmount:  row.TotalAmount,
	}
}

func toOrderInfo(row *model.OrderRow, items []*model.OrderItemRow) *order.OrderInfo {
	if row == nil {
		return nil
	}
	info := &order.OrderInfo{
		Id:                 row.Id,
		OrderNo:            row.OrderNo,
		UserId:             row.UserId,
		TotalAmount:        row.TotalAmount,
		FreightAmount:      row.FreightAmount,
		DiscountAmount:     row.DiscountAmount,
		PayAmount:          row.PayAmount,
		Status:             int32(row.Status),
		PayType:            int32(row.PayType),
		TransactionId:      row.TransactionId,
		ReceiverName:       row.ReceiverName,
		ReceiverPhone:      row.ReceiverPhone,
		ReceiverProvince:   row.ReceiverProvince,
		ReceiverCity:       row.ReceiverCity,
		ReceiverDistrict:   row.ReceiverDistrict,
		ReceiverAddress:    row.ReceiverAddress,
		ReceiverPostalCode: row.ReceiverPostalCode,
		Remark:             row.Remark,
		AdminRemark:        row.AdminRemark,
		LogisticsCompany:   row.LogisticsCompany,
		TrackingNo:         row.TrackingNo,
		PaidAt:             row.PaidAt,
		ShippedAt:          row.ShippedAt,
		CompletedAt:        row.CompletedAt,
		CanceledAt:         row.CanceledAt,
		CancelReason:       row.CancelReason,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
	for _, item := range items {
		info.Items = append(info.Items, toOrderItemInfo(item))
	}
	return info
}

func toCartItemInfo(row *model.CartItemRow) *order.CartItemInfo {
	return &order.CartItemInfo{
		Id:            row.Id,
		ProductId:     row.ProductId,
		ProductName:   row.ProductName,
		ProductCover:  row.ProductCover,
		Price:         row.Price,
		Stock:         int32(row.Stock),
		Quantity:      int32(row.Quantity),
		Selected:      int32(row.Selected),
		ProductStatus: int32(row.ProductStatus),
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func toAfterSaleInfo(row *model.AfterSaleRow) *order.AfterSaleInfo {
	if row == nil {
		return nil
	}
	return &order.AfterSaleInfo{
		Id:               row.Id,
		AfterSaleNo:      row.AfterSaleNo,
		OrderId:          row.OrderId,
		OrderItemId:      row.OrderItemId,
		UserId:           row.UserId,
		Type:             int32(row.Type),
		Reason:           row.Reason,
		Description:      row.Description,
		Images:           decodeImages(row.Images),
		Amount:           row.Amount,
		Status:           int32(row.Status),
		AdminRemark:      row.AdminRemark,
		TrackingNo:       row.TrackingNo,
		ApplyAt:          row.ApplyAt,
		ReviewAt:         row.ReviewAt,
		BuyerShippedAt:   row.BuyerShippedAt,
		SellerReceivedAt: row.SellerReceivedAt,
		RefundedAt:       row.RefundedAt,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
		OrderNo:          row.OrderNo,
		ProductName:      row.ProductName,
		ProductCover:     row.ProductCover,
	}
}

// loadOrderInfo loads an order and its items inside or outside a transaction.
func loadOrderInfo(ctx context.Context, svcCtx *svc.ServiceContext, session sqlx.Session, id, userID int64, admin bool) (*order.OrderInfo, error) {
	row, err := svcCtx.OrdersModel.FindOrder(ctx, session, id, userID, admin)
	if err != nil {
		return nil, err
	}
	items, err := svcCtx.OrderItemModel.FindOrderItems(ctx, session, id)
	if err != nil {
		return nil, err
	}
	return toOrderInfo(row, items), nil
}
