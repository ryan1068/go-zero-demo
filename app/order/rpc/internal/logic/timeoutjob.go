package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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
func StartTimeoutJob(db *sql.DB) {
	go func() {
		ticker := time.NewTicker(timeoutScanInterval)
		defer ticker.Stop()
		for range ticker.C {
			if err := cancelExpiredOrders(db); err != nil {
				logx.Errorf("cancel expired orders failed: %v", err)
			}
			if err := completeExpiredShippedOrders(db); err != nil {
				logx.Errorf("auto confirm shipped orders failed: %v", err)
			}
		}
	}()
}

func cancelExpiredOrders(db *sql.DB) error {
	ctx := context.Background()
	rows, err := db.QueryContext(ctx,
		"SELECT id FROM orders WHERE status = ? AND deleted_at IS NULL AND created_at < DATE_SUB(NOW(), INTERVAL ? MINUTE) ORDER BY id LIMIT ?",
		orderStatusPendingPayment, timeoutOrderMinutes, timeoutBatchSize,
	)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logx.Errorf("close expired order rows failed: %v", closeErr)
		}
	}()

	ids := make([]int64, 0, timeoutBatchSize)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, id := range ids {
		if err := cancelExpiredOrder(ctx, db, id); err != nil {
			logx.Errorf("cancel expired order %d failed: %v", id, err)
		}
	}
	return nil
}

func cancelExpiredOrder(ctx context.Context, db *sql.DB, id int64) error {
	logger := logx.WithContext(ctx)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollbackTx(logger, tx)

	var status int32
	err = tx.QueryRowContext(ctx,
		"SELECT status FROM orders WHERE id = ? AND deleted_at IS NULL FOR UPDATE", id,
	).Scan(&status)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return err
	}
	if status != orderStatusPendingPayment {
		return tx.Commit()
	}

	if _, err := tx.ExecContext(ctx,
		"UPDATE orders SET status = ?, canceled_at = NOW(), cancel_reason = ? WHERE id = ?",
		orderStatusCanceled, "超时未支付", id,
	); err != nil {
		return err
	}
	if err := restoreStock(ctx, tx, id); err != nil {
		return err
	}
	if err := writeOrderLog(ctx, tx, id, orderStatusPendingPayment, orderStatusCanceled, operatorSystem, 0, "超时未支付"); err != nil {
		return err
	}
	return tx.Commit()
}

func completeExpiredShippedOrders(db *sql.DB) error {
	ctx := context.Background()
	rows, err := db.QueryContext(ctx,
		"SELECT id FROM orders WHERE status = ? AND deleted_at IS NULL AND shipped_at IS NOT NULL AND shipped_at < DATE_SUB(NOW(), INTERVAL ? DAY) ORDER BY id LIMIT ?",
		orderStatusShipped, autoConfirmDays, autoConfirmBatchSize,
	)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logx.Errorf("close auto confirm order rows failed: %v", closeErr)
		}
	}()

	ids := make([]int64, 0, autoConfirmBatchSize)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, id := range ids {
		if err := completeExpiredShippedOrder(ctx, db, id); err != nil {
			logx.Errorf("auto confirm order %d failed: %v", id, err)
		}
	}
	return nil
}

func completeExpiredShippedOrder(ctx context.Context, db *sql.DB, id int64) error {
	logger := logx.WithContext(ctx)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollbackTx(logger, tx)

	var status int32
	err = tx.QueryRowContext(ctx,
		"SELECT status FROM orders WHERE id = ? AND deleted_at IS NULL FOR UPDATE", id,
	).Scan(&status)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return err
	}
	if status != orderStatusShipped {
		return tx.Commit()
	}

	if _, err := tx.ExecContext(ctx,
		"UPDATE orders SET status = ?, completed_at = NOW() WHERE id = ? AND status = ?",
		orderStatusCompleted, id, orderStatusShipped,
	); err != nil {
		return err
	}
	if err := writeOrderLog(ctx, tx, id, orderStatusShipped, orderStatusCompleted, operatorSystem, 0, "超时自动确认收货"); err != nil {
		return err
	}
	return tx.Commit()
}
