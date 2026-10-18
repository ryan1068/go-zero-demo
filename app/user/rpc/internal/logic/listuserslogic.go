package logic

import (
	"context"
	"strings"
	"unicode/utf8"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/model"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	defaultUserPageSize = 10
	maxUserPageSize     = 100
)

type ListUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUsersLogic) ListUsers(in *user.ListUsersReq) (*user.ListUsersResp, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = defaultUserPageSize
	}
	if pageSize > maxUserPageSize {
		pageSize = maxUserPageSize
	}

	username := strings.TrimSpace(in.Username)
	nickname := strings.TrimSpace(in.Nickname)
	email := strings.TrimSpace(in.Email)
	phone := strings.TrimSpace(in.Phone)
	if utf8.RuneCountInString(username) > maxUsernameLength ||
		utf8.RuneCountInString(nickname) > maxNicknameLength ||
		utf8.RuneCountInString(email) > maxEmailLength ||
		utf8.RuneCountInString(phone) > maxPhoneLength {
		return &user.ListUsersResp{
			Code: errcode.BadRequest,
			Msg:  errcode.Msg(errcode.BadRequest),
		}, nil
	}

	list, total, err := l.svcCtx.UserModel.FindUsers(l.ctx, model.UserFilter{
		Username: username,
		Nickname: nickname,
		Email:    email,
		Phone:    phone,
		Status:   int64(in.Status),
		Page:     int64(page),
		PageSize: int64(pageSize),
	})
	if err != nil {
		l.Errorf("query users failed: %v", err)
		return nil, err
	}

	items := make([]*user.UserInfo, 0, len(list))
	for _, item := range list {
		items = append(items, toUserInfo(item))
	}

	return &user.ListUsersResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
