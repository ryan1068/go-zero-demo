package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRoleLogic {
	return &UpdateUserRoleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserRoleLogic) UpdateUserRole(req *types.UpdateUserRoleReq) (*types.BaseResp, error) {
	ctx, _ := injectAdmin(l.ctx)
	s, err := l.svcCtx.UserRpc.UpdateUserRole(ctx, &user.UpdateUserRoleReq{Id: req.Id, Role: req.Role})
	if err != nil {
		l.Errorf("update user role rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}

type ListAuditLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAuditLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAuditLogsLogic {
	return &ListAuditLogsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListAuditLogsLogic) ListAuditLogs(req *types.ListAuditLogsReq) (*types.ListAuditLogsResp, error) {
	s, err := l.svcCtx.UserRpc.ListAuditLogs(l.ctx, &user.ListAuditLogsReq{
		Page:       req.Page,
		PageSize:   req.PageSize,
		AdminId:    req.AdminId,
		Action:     req.Action,
		TargetType: req.TargetType,
	})
	if err != nil {
		l.Errorf("list audit logs rpc failed: %v", err)
		return nil, err
	}
	list := make([]types.AuditLogInfo, 0, len(s.List))
	for _, item := range s.List {
		list = append(list, types.AuditLogInfo{
			Id:         item.Id,
			AdminId:    item.AdminId,
			AdminName:  item.AdminName,
			Action:     item.Action,
			TargetType: item.TargetType,
			TargetId:   item.TargetId,
			Detail:     item.Detail,
			CreatedAt:  item.CreatedAt,
		})
	}
	return &types.ListAuditLogsResp{
		Code: s.Code,
		Msg:  s.Msg,
		Data: types.AuditLogListData{
			List:     list,
			Total:    s.Total,
			Page:     s.Page,
			PageSize: s.PageSize,
		},
	}, nil
}
