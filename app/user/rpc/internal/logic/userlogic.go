package logic

import (
	"context"
	"go-zero-demo/user-rpc/pb/user"

	"go-zero-demo/user-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
	return &user.UserResp{
		Message: in.Name + "345",
	}, nil
}
