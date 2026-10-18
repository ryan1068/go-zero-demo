package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	username := strings.TrimSpace(in.Username)
	if username == "" || in.Password == "" {
		return &user.LoginResp{
			Code: errcode.BadRequest,
			Msg:  errcode.Msg(errcode.BadRequest),
		}, nil
	}

	stored, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, username)
	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		return &user.LoginResp{
			Code: errcode.Unauthorized,
			Msg:  errcode.Msg(errcode.Unauthorized),
		}, nil
	case err != nil:
		l.Errorf("query user by username failed: %v", err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(in.Password)); err != nil {
		return &user.LoginResp{
			Code: errcode.Unauthorized,
			Msg:  errcode.Msg(errcode.Unauthorized),
		}, nil
	}
	if stored.Status != 1 {
		return &user.LoginResp{
			Code: errcode.Forbidden,
			Msg:  errcode.Msg(errcode.Forbidden),
		}, nil
	}

	ttl := time.Duration(l.svcCtx.Config.Jwt.ExpireHours) * time.Hour
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	token, err := auth.SignWithRole(l.svcCtx.Config.Jwt.Secret, stored.Id, stored.Username, int32(stored.Role), ttl)
	if err != nil {
		l.Errorf("sign token failed: %v", err)
		return nil, err
	}
	if err := l.svcCtx.UserModel.UpdateLastLogin(l.ctx, stored.Id); err != nil {
		l.Errorf("update last login time failed: %v", err)
	}

	return &user.LoginResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		Token:    token,
		UserId:   stored.Id,
		Username: stored.Username,
		Role:     int32(stored.Role),
	}, nil
}
