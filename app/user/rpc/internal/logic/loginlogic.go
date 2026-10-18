package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
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
	username := strings.TrimSpace(in.Username)
	if username == "" || in.Password == "" {
		return &user.LoginResp{
			Code: errcode.BadRequest,
			Msg:  errcode.Msg(errcode.BadRequest),
		}, nil
	}

	var (
		id       int64
		name     string
		password string
		status   int32
		role     int32
	)
	err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT id, username, password, status, role FROM user WHERE username = ?", username,
	).Scan(&id, &name, &password, &status, &role)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &user.LoginResp{
			Code: errcode.Unauthorized,
			Msg:  errcode.Msg(errcode.Unauthorized),
		}, nil
	case err != nil:
		l.Errorf("query user by username failed: %v", err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(in.Password)); err != nil {
		return &user.LoginResp{
			Code: errcode.Unauthorized,
			Msg:  errcode.Msg(errcode.Unauthorized),
		}, nil
	}
	if status != 1 {
		return &user.LoginResp{
			Code: errcode.Forbidden,
			Msg:  errcode.Msg(errcode.Forbidden),
		}, nil
	}

	ttl := time.Duration(l.svcCtx.Config.Jwt.ExpireHours) * time.Hour
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	token, err := auth.SignWithRole(l.svcCtx.Config.Jwt.Secret, id, name, role, ttl)
	if err != nil {
		l.Errorf("sign token failed: %v", err)
		return nil, err
	}
	if _, err := l.svcCtx.DB.ExecContext(l.ctx,
		"UPDATE user SET last_login_at = NOW() WHERE id = ?", id,
	); err != nil {
		l.Errorf("update last login time failed: %v", err)
	}

	return &user.LoginResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		Token:    token,
		UserId:   id,
		Username: name,
		Role:     role,
	}, nil
}
