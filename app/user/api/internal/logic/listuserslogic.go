package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.ListUsersReq) (*types.ListUsersResp, error) {
	s, err := l.svcCtx.UserRpc.ListUsers(l.ctx, &user.ListUsersReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	})
	if err != nil {
		l.Errorf("list users rpc failed: %v", err)
		return nil, err
	}

	resp := &types.ListUsersResp{
		Code: s.Code,
		Msg:  s.Msg,
		Data: types.UserListData{
			List:     make([]types.UserInfo, 0, len(s.List)),
			Total:    s.Total,
			Page:     s.Page,
			PageSize: s.PageSize,
		},
	}
	for _, item := range s.List {
		if item == nil {
			continue
		}
		resp.Data.List = append(resp.Data.List, toUserInfo(item))
	}
	return resp, nil
}
