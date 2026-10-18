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

type UpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserStatusLogic) UpdateUserStatus(in *user.UpdateUserStatusReq) (*user.BaseResp, error) {
	if in.Id <= 0 || (in.Status != 0 && in.Status != 1) {
		return &user.BaseResp{
			Code: errcode.BadRequest,
			Msg:  errcode.Msg(errcode.BadRequest),
		}, nil
	}

	affected, err := l.svcCtx.UserModel.UpdateStatus(l.ctx, in.Id, int64(in.Status))
	if err != nil {
		l.Errorf("update user status failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		_, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
		if errors.Is(err, sqlx.ErrNotFound) {
			return &user.BaseResp{
				Code: errcode.NotFound,
				Msg:  errcode.Msg(errcode.NotFound),
			}, nil
		}
		if err != nil {
			l.Errorf("query user failed: %v", err)
			return nil, err
		}
	}

	adminID, _ := auth.UserIDFromIncoming(l.ctx)
	if err := l.svcCtx.AdminAuditLogModel.InsertAudit(l.ctx, adminID, "update_user_status", "user", in.Id, fmt.Sprintf("status=%d", in.Status)); err != nil {
		l.Errorf("write audit log failed: %v", err)
	}
	return &user.BaseResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
	}, nil
}
