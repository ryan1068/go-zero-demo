package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
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

	result, err := l.svcCtx.DB.ExecContext(l.ctx,
		"UPDATE user SET status = ? WHERE id = ?", in.Status, in.Id,
	)
	if err != nil {
		l.Errorf("update user status failed: %v", err)
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		l.Errorf("read affected rows failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		if _, err := findUserByID(l.ctx, l.svcCtx.DB, in.Id); errors.Is(err, sql.ErrNoRows) {
			return &user.BaseResp{
				Code: errcode.NotFound,
				Msg:  errcode.Msg(errcode.NotFound),
			}, nil
		} else if err != nil {
			l.Errorf("query user failed: %v", err)
			return nil, err
		}
	}
	adminID, _ := auth.UserIDFromIncoming(l.ctx)
	auth.Audit(l.ctx, l.svcCtx.DB, adminID, "update_user_status", "user", in.Id, fmt.Sprintf("status=%d", in.Status))
	return &user.BaseResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
	}, nil
}
