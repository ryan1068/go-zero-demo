package logic

import (
	"context"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

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
	info, code, err := createUser(l.ctx, l.svcCtx, createUserParams{
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
	if err := l.svcCtx.AdminAuditLogModel.InsertAudit(l.ctx, adminID, "create_user", "user", info.Id, in.Username); err != nil {
		l.Errorf("write audit log failed: %v", err)
	}

	return &user.CreateUserResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		Data: info,
	}, nil
}
