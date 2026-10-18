package logic

import (
	"context"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	info, code, err := createUser(l.ctx, l.svcCtx, createUserParams{
		Username: in.Username,
		Password: in.Password,
		Nickname: in.Nickname,
		Email:    in.Email,
		Phone:    in.Phone,
		Status:   1,
		Role:     0,
	})
	if err != nil {
		l.Errorf("register user failed: %v", err)
		return nil, err
	}
	if code != errcode.Success {
		return &user.RegisterResp{
			Code: code,
			Msg:  errcode.Msg(code),
		}, nil
	}

	return &user.RegisterResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		Data: info,
	}, nil
}
