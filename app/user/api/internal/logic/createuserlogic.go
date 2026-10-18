package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (*types.CreateUserResp, error) {
	ctx, _ := injectAdmin(l.ctx)
	s, err := l.svcCtx.UserRpc.CreateUser(ctx, &user.CreateUserReq{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Remark:   req.Remark,
		Role:     req.Role,
	})
	if err != nil {
		l.Errorf("create user rpc failed: %v", err)
		return nil, err
	}

	resp := &types.CreateUserResp{
		Code: s.Code,
		Msg:  s.Msg,
	}
	if s.Data != nil {
		resp.Data = toUserInfo(s.Data)
	}
	return resp, nil
}
