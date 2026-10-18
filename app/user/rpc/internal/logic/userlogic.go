package logic

import (
	"context"
	"errors"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserLogic) User(in *user.UserReq) (*user.UserResp, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return &user.UserResp{
			Code: errcode.BadRequest,
			Msg:  errcode.Msg(errcode.BadRequest),
		}, nil
	}

	stored, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, name)
	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		return &user.UserResp{
			Code: errcode.NotFound,
			Msg:  errcode.Msg(errcode.NotFound),
		}, nil
	case err != nil:
		l.Errorf("query user failed: %v", err)
		return nil, err
	}

	return &user.UserResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		Data: toUserInfo(stored),
	}, nil
}
