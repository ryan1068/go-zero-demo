package auth

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

// Audit records a management operation in the shared audit table. Failures are
// logged but never block the business operation.
func Audit(ctx context.Context, db *sql.DB, adminID int64, action, targetType string, targetID int64, detail string) {
	if db == nil {
		return
	}
	if _, err := db.ExecContext(ctx,
		"INSERT INTO looklook_usercenter.admin_audit_log (admin_id, action, target_type, target_id, detail) VALUES (?, ?, ?, ?, ?)",
		adminID, action, targetType, targetID, detail,
	); err != nil {
		logx.WithContext(ctx).Errorf("write admin audit log failed: %v", err)
	}
}
