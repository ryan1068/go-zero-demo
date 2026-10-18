package logic

import (
	"context"
	"strings"
	"unicode/utf8"

	"go-zero-demo/pkg/errcode"
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

	where := []string{"1 = 1"}
	args := make([]any, 0, 5)
	if username != "" {
		where = append(where, "username LIKE ?")
		args = append(args, username+"%")
	}
	if nickname != "" {
		where = append(where, "nickname LIKE ?")
		args = append(args, nickname+"%")
	}
	if email != "" {
		where = append(where, "email LIKE ?")
		args = append(args, email+"%")
	}
	if phone != "" {
		where = append(where, "phone LIKE ?")
		args = append(args, phone+"%")
	}
	// status: 0 = all, 1 = active, 2 = disabled
	switch in.Status {
	case 1:
		where = append(where, "status = ?")
		args = append(args, 1)
	case 2:
		where = append(where, "status = ?")
		args = append(args, 0)
	}
	condition := strings.Join(where, " AND ")

	var total int64
	countQuery := "SELECT COUNT(*) FROM user WHERE " + condition
	if err := l.svcCtx.DB.QueryRowContext(l.ctx, countQuery, args...).Scan(&total); err != nil {
		l.Errorf("count users failed: %v", err)
		return nil, err
	}

	offset := int64(page-1) * int64(pageSize)
	listArgs := append(append([]any{}, args...), int64(pageSize), offset)
	rows, err := l.svcCtx.DB.QueryContext(l.ctx,
		"SELECT "+userSelectColumns+" FROM user WHERE "+condition+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	if err != nil {
		l.Errorf("query users failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close user rows failed: %v", closeErr)
		}
	}()

	list := make([]*user.UserInfo, 0, pageSize)
	for rows.Next() {
		info, err := scanUserInfo(rows)
		if err != nil {
			l.Errorf("scan user failed: %v", err)
			return nil, err
		}
		list = append(list, info)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate users failed: %v", err)
		return nil, err
	}

	return &user.ListUsersResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
