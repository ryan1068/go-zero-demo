package logic

import (
	"context"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAuditLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAuditLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAuditLogsLogic {
	return &ListAuditLogsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListAuditLogsLogic) ListAuditLogs(in *user.ListAuditLogsReq) (*user.ListAuditLogsResp, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	where := []string{"1 = 1"}
	args := make([]any, 0, 3)
	if in.AdminId > 0 {
		where = append(where, "a.admin_id = ?")
		args = append(args, in.AdminId)
	}
	if action := strings.TrimSpace(in.Action); action != "" {
		where = append(where, "a.action = ?")
		args = append(args, action)
	}
	if targetType := strings.TrimSpace(in.TargetType); targetType != "" {
		where = append(where, "a.target_type = ?")
		args = append(args, targetType)
	}
	condition := strings.Join(where, " AND ")

	var total int64
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM admin_audit_log a WHERE "+condition, args...,
	).Scan(&total); err != nil {
		l.Errorf("count audit logs failed: %v", err)
		return nil, err
	}

	offset := int64(page-1) * int64(pageSize)
	listArgs := append(append([]any{}, args...), int64(pageSize), offset)
	rows, err := l.svcCtx.DB.QueryContext(l.ctx,
		"SELECT a.id, a.admin_id, COALESCE(u.username, ''), a.action, a.target_type, a.target_id, a.detail, UNIX_TIMESTAMP(a.created_at)"+
			" FROM admin_audit_log a LEFT JOIN user u ON u.id = a.admin_id"+
			" WHERE "+condition+" ORDER BY a.id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	if err != nil {
		l.Errorf("query audit logs failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close audit log rows failed: %v", closeErr)
		}
	}()

	list := make([]*user.AuditLogInfo, 0, pageSize)
	for rows.Next() {
		item := new(user.AuditLogInfo)
		if err := rows.Scan(&item.Id, &item.AdminId, &item.AdminName, &item.Action, &item.TargetType, &item.TargetId, &item.Detail, &item.CreatedAt); err != nil {
			l.Errorf("scan audit log failed: %v", err)
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate audit logs failed: %v", err)
		return nil, err
	}
	return &user.ListAuditLogsResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
