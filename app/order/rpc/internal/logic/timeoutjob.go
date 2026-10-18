package logic

import (
	"context"
	"errors"
	"time"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	timeoutScanInterval  = time.Minute
	timeoutOrderMinutes  = 30
	timeoutBatchSize     = 100
	autoConfirmDays      = 7
	autoConfirmBatchSize = 100
)

// StartTimeoutJob periodically cancels unpaid orders that have exceeded the
// payment window, restores their stock, and auto-confirms shipped orders.
func StartTimeoutJob(svcCtx *svc.ServiceContext) {
	go func() {
		ticker := time.NewTicker(timeoutScanInterval)
		defer ticker.Stop()
		for range ticker.C {
			if err := cancelExpiredOrders(svcCtx); err != nil {
				logx.Errorf("cancel expired orders failed: %v", err)
			}
			if err := completeExpiredShippedOrders(svcCtx); err != nil {
				logx.Errorf("auto confirm shipped orders failed: %v", err)
			}
		}
	}()
}

func cancelExpiredOrders(svcCtx *svc.ServiceContext) error {
	ctx := context.Background()
	ids, err := svcCtx.OrdersModel.FindExpiredPendingOrderIDs(ctx, nil,
		orderStatusPendingPayment, timeoutOrderMinutes, timeoutBatchSize)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := cancelExpiredOrder(ctx, svcCtx, id); err != nil {
			logx.Errorf("cancel expired order %d failed: %v", id, err)
		}
	}
	return nil
}

func cancelExpiredOrder(ctx context.Context, svcCtx *svc.ServiceContext, id int64) error {
	err := svcCtx.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		lock, err := svcCtx.OrdersModel.LockOrder(ctx, session, id, nil)
		if errors.Is(err, sqlx.ErrNotFound) {
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}
		if lock.Status != orderStatusPendingPayment {
			// Preserve the original behavior: wrong status commits a no-op.
			return nil
		}

		if _, err := svcCtx.OrdersModel.UpdateOrderCanceled(ctx, session, id, orderStatusCanceled, "超时未支付"); err != nil {
			return err
		}
		lines, err := svcCtx.OrderItemModel.StockLines(ctx, session, id)
		if err != nil {
			return err
		}
		for _, line := range lines {
			if err := svcCtx.ProductModel.RestoreStock(ctx, session, line.ProductId, line.Quantity); err != nil {
				return err
			}
		}
		return svcCtx.OrderLogModel.InsertOrderLog(ctx, session, model.NewOrderLog{
			OrderID:      id,
			FromStatus:   orderStatusPendingPayment,
			ToStatus:     orderStatusCanceled,
			OperatorType: operatorSystem,
			OperatorID:   0,
			Remark:       "超时未支付",
		})
	})
	if errors.Is(err, model.ErrAbortTransaction) {
		return nil
	}
	return err
}

func completeExpiredShippedOrders(svcCtx *svc.ServiceContext) error {
	ctx := context.Background()
	ids, err := svcCtx.OrdersModel.FindExpiredShippedOrderIDs(ctx, nil,
		orderStatusShipped, autoConfirmDays, autoConfirmBatchSize)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := completeExpiredShippedOrder(ctx, svcCtx, id); err != nil {
			logx.Errorf("auto confirm order %d failed: %v", id, err)
		}
	}
	return nil
}

func completeExpiredShippedOrder(ctx context.Context, svcCtx *svc.ServiceContext, id int64) error {
	err := svcCtx.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		lock, err := svcCtx.OrdersModel.LockOrder(ctx, session, id, nil)
		if errors.Is(err, sqlx.ErrNotFound) {
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}
		if lock.Status != orderStatusShipped {
			// Preserve the original behavior: wrong status commits a no-op.
			return nil
		}

		if _, err := svcCtx.OrdersModel.UpdateOrderCompletedByStatus(ctx, session, id, orderStatusShipped, orderStatusCompleted); err != nil {
			return err
		}
		return svcCtx.OrderLogModel.InsertOrderLog(ctx, session, model.NewOrderLog{
			OrderID:      id,
			FromStatus:   orderStatusShipped,
			ToStatus:     orderStatusCompleted,
			OperatorType: operatorSystem,
			OperatorID:   0,
			Remark:       "超时自动确认收货",
		})
	})
	if errors.Is(err, model.ErrAbortTransaction) {
		return nil
	}
	return err
}
