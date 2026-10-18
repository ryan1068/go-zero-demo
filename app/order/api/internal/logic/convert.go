package logic

import (
	"context"

	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
)

func injectUser(ctx context.Context) (context.Context, bool) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return ctx, false
	}
	return auth.InjectUserID(ctx, userID), true
}

func toCartItem(info *order.CartItemInfo) types.CartItemInfo {
	return types.CartItemInfo{
		Id:            info.Id,
		ProductId:     info.ProductId,
		ProductName:   info.ProductName,
		ProductCover:  info.ProductCover,
		Price:         info.Price,
		Stock:         info.Stock,
		Quantity:      info.Quantity,
		Selected:      info.Selected,
		ProductStatus: info.ProductStatus,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}
}

func toOrderInfo(info *order.OrderInfo) types.OrderInfo {
	result := types.OrderInfo{
		Id:                 info.Id,
		OrderNo:            info.OrderNo,
		UserId:             info.UserId,
		TotalAmount:        info.TotalAmount,
		FreightAmount:      info.FreightAmount,
		DiscountAmount:     info.DiscountAmount,
		PayAmount:          info.PayAmount,
		Status:             info.Status,
		PayType:            info.PayType,
		TransactionId:      info.TransactionId,
		ReceiverName:       info.ReceiverName,
		ReceiverPhone:      info.ReceiverPhone,
		ReceiverProvince:   info.ReceiverProvince,
		ReceiverCity:       info.ReceiverCity,
		ReceiverDistrict:   info.ReceiverDistrict,
		ReceiverAddress:    info.ReceiverAddress,
		ReceiverPostalCode: info.ReceiverPostalCode,
		Remark:             info.Remark,
		AdminRemark:        info.AdminRemark,
		LogisticsCompany:   info.LogisticsCompany,
		TrackingNo:         info.TrackingNo,
		PaidAt:             info.PaidAt,
		ShippedAt:          info.ShippedAt,
		CompletedAt:        info.CompletedAt,
		CanceledAt:         info.CanceledAt,
		CancelReason:       info.CancelReason,
		CreatedAt:          info.CreatedAt,
		UpdatedAt:          info.UpdatedAt,
		Items:              make([]types.OrderItemInfo, 0, len(info.Items)),
	}
	for _, item := range info.Items {
		if item == nil {
			continue
		}
		result.Items = append(result.Items, types.OrderItemInfo{
			Id:           item.Id,
			ProductId:    item.ProductId,
			ProductName:  item.ProductName,
			ProductCover: item.ProductCover,
			Price:        item.Price,
			Quantity:     item.Quantity,
			TotalAmount:  item.TotalAmount,
		})
	}
	return result
}

func toAfterSaleInfo(info *order.AfterSaleInfo) types.AfterSaleInfo {
	return types.AfterSaleInfo{
		Id:               info.Id,
		AfterSaleNo:      info.AfterSaleNo,
		OrderId:          info.OrderId,
		OrderItemId:      info.OrderItemId,
		UserId:           info.UserId,
		Type:             info.Type,
		Reason:           info.Reason,
		Description:      info.Description,
		Images:           info.Images,
		Amount:           info.Amount,
		Status:           info.Status,
		AdminRemark:      info.AdminRemark,
		TrackingNo:       info.TrackingNo,
		ApplyAt:          info.ApplyAt,
		ReviewAt:         info.ReviewAt,
		BuyerShippedAt:   info.BuyerShippedAt,
		SellerReceivedAt: info.SellerReceivedAt,
		RefundedAt:       info.RefundedAt,
		CreatedAt:        info.CreatedAt,
		UpdatedAt:        info.UpdatedAt,
		OrderNo:          info.OrderNo,
		ProductName:      info.ProductName,
		ProductCover:     info.ProductCover,
	}
}
