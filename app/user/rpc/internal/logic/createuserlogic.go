package logic

import (
	"context"
	"regexp"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 64
	minPasswordLength = 6
	maxPasswordLength = 72
	maxNicknameLength = 255
	maxEmailLength    = 128
	maxPhoneLength    = 32
	maxAvatarLength   = 512
	maxRemarkLength   = 255
)

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *user.CreateUserReq) (*user.CreateUserResp, error) {
	info, code, err := createUser(l.ctx, l.svcCtx.DB, createUserParams{
		Username: in.Username,
		Password: in.Password,
		Nickname: in.Nickname,
		Email:    in.Email,
		Phone:    in.Phone,
		Avatar:   in.Avatar,
		Status:   in.Status,
		Remark:   in.Remark,
		Role:     in.Role,
	})
	if err != nil {
		l.Errorf("create user failed: %v", err)
		return nil, err
	}
	if code != errcode.Success {
		return &user.CreateUserResp{
			Code: code,
			Msg:  errcode.Msg(code),
		}, nil
	}

	adminID, _ := auth.UserIDFromIncoming(l.ctx)
	auth.Audit(l.ctx, l.svcCtx.DB, adminID, "create_user", "user", info.Id, in.Username)

	return &user.CreateUserResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		Data: info,
	}, nil
}
