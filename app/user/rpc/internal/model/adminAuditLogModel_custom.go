package model

import (
	"context"
	"strings"
)

// AuditLogFilter is the filter for the admin audit log list.
type AuditLogFilter struct {
	AdminID    int64
	Action     string
	TargetType string
	Page       int64
	PageSize   int64
}

// AuditLogRow is an audit log joined with the operator name.
type AuditLogRow struct {
	Id         int64  `db:"id"`
	AdminId    int64  `db:"admin_id"`
	AdminName  string `db:"admin_name"`
	Action     string `db:"action"`
	TargetType string `db:"target_type"`
	TargetId   int64  `db:"target_id"`
	Detail     string `db:"detail"`
	CreatedAt  int64  `db:"created_at"`
}

// InsertAudit writes one admin audit record.
func (m *customAdminAuditLogModel) InsertAudit(ctx context.Context, adminID int64, action, targetType string, targetID int64, detail string) error {
	_, err := m.Insert(ctx, &AdminAuditLog{
		AdminId:    adminID,
		Action:     action,
		TargetType: targetType,
		TargetId:   targetID,
		Detail:     detail,
	})
	return err
}

// ListAuditLogs returns a page of audit logs and the total count.
func (m *customAdminAuditLogModel) ListAuditLogs(ctx context.Context, f AuditLogFilter) ([]*AuditLogRow, int64, error) {
	where := []string{"1 = 1"}
	args := make([]any, 0, 3)
	if f.AdminID > 0 {
		where = append(where, "a.`admin_id` = ?")
		args = append(args, f.AdminID)
	}
	if f.Action != "" {
		where = append(where, "a.`action` = ?")
		args = append(args, f.Action)
	}
	if f.TargetType != "" {
		where = append(where, "a.`target_type` = ?")
		args = append(args, f.TargetType)
	}
	condition := strings.Join(where, " AND ")

	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total,
		"select count(*) from "+m.table+" a where "+condition, args...,
	); err != nil {
		return nil, 0, err
	}

	listArgs := append(append([]any{}, args...), f.PageSize, (f.Page-1)*f.PageSize)
	var list []*AuditLogRow
	if err := m.conn.QueryRowsCtx(ctx, &list,
		"select a.`id`, a.`admin_id`, coalesce(u.`username`, '') as admin_name, a.`action`,"+
			" a.`target_type`, a.`target_id`, a.`detail`, unix_timestamp(a.`created_at`) as created_at"+
			" from "+m.table+" a left join `user` u on u.`id` = a.`admin_id`"+
			" where "+condition+" order by a.`id` desc limit ? offset ?",
		listArgs...,
	); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
