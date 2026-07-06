package logic

import (
	"context"
	"fmt"
	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) User(req *types.UserReq) (resp *types.UserResp, err error) {
	s, err := l.svcCtx.UserRpc.User(l.ctx, &user.UserReq{
		Name: req.Name,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp = new(types.UserResp)
	resp.Message = s.Message
	return
}
