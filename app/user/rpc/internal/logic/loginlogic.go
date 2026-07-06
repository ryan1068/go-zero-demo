package logic

import (
	"context"

	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	var (
		id       int64
		username string
		password string
	)
	err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT id, username, password FROM user WHERE username = ?", in.Username,
	).Scan(&id, &username, &password)

	if err != nil {
		return &user.LoginResp{
			Code: 401,
			Msg:  "用户名或密码错误",
		}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(in.Password)); err != nil {
		return &user.LoginResp{
			Code: 401,
			Msg:  "用户名或密码错误",
		}, nil
	}

	return &user.LoginResp{
		Code: 200,
		Msg:  "登录成功",
	}, nil
}
