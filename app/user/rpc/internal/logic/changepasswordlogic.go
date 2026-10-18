package logic

import (
	"context"
	"errors"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ChangePasswordLogic) ChangePassword(in *user.ChangePasswordReq) (*user.BaseResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &user.BaseResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if len(in.OldPassword) < minPasswordLength || len(in.OldPassword) > maxPasswordLength ||
		len(in.NewPassword) < minPasswordLength || len(in.NewPassword) > maxPasswordLength {
		return &user.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	if in.OldPassword == in.NewPassword {
		return &user.BaseResp{Code: errcode.BadRequest, Msg: "新密码不能与原密码相同"}, nil
	}

	stored, err := l.svcCtx.UserModel.FindOne(l.ctx, userID)
	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		return &user.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query password failed: %v", err)
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(in.OldPassword)); err != nil {
		return &user.BaseResp{Code: errcode.Unauthorized, Msg: "原密码错误"}, nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("hash new password failed: %v", err)
		return nil, err
	}
	if err := l.svcCtx.UserModel.UpdatePassword(l.ctx, userID, string(hash)); err != nil {
		l.Errorf("update password failed: %v", err)
		return nil, err
	}
	return &user.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
