package logic

import (
	"context"
	"errors"
	"fmt"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRoleLogic {
	return &UpdateUserRoleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateUserRoleLogic) UpdateUserRole(in *user.UpdateUserRoleReq) (*user.BaseResp, error) {
	if in.Id <= 0 || (in.Role != 0 && in.Role != 1 && in.Role != 2) {
		return &user.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	affected, err := l.svcCtx.UserModel.UpdateRole(l.ctx, in.Id, int64(in.Role))
	if err != nil {
		l.Errorf("update user role failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		_, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
		if errors.Is(err, sqlx.ErrNotFound) {
			return &user.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
		}
		if err != nil {
			l.Errorf("query user failed: %v", err)
			return nil, err
		}
	}

	adminID, _ := auth.UserIDFromIncoming(l.ctx)
	if err := l.svcCtx.AdminAuditLogModel.InsertAudit(l.ctx, adminID, "update_user_role", "user", in.Id, fmt.Sprintf("role=%d", in.Role)); err != nil {
		l.Errorf("write audit log failed: %v", err)
	}
	return &user.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
