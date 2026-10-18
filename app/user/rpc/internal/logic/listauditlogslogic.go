package logic

import (
	"context"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/model"
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

	rows, total, err := l.svcCtx.AdminAuditLogModel.ListAuditLogs(l.ctx, model.AuditLogFilter{
		AdminID:    in.AdminId,
		Action:     strings.TrimSpace(in.Action),
		TargetType: strings.TrimSpace(in.TargetType),
		Page:       int64(page),
		PageSize:   int64(pageSize),
	})
	if err != nil {
		l.Errorf("query audit logs failed: %v", err)
		return nil, err
	}

	list := make([]*user.AuditLogInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, &user.AuditLogInfo{
			Id:         row.Id,
			AdminId:    row.AdminId,
			Action:     row.Action,
			TargetType: row.TargetType,
			TargetId:   row.TargetId,
			Detail:     row.Detail,
			CreatedAt:  row.CreatedAt,
			AdminName:  row.AdminName,
		})
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
