package logic

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/model"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"
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
	mysqlDuplicateKey = 1062
)

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type createUserParams struct {
	Username string
	Password string
	Nickname string
	Email    string
	Phone    string
	Avatar   string
	Status   int32
	Remark   string
	Role     int32
}

// createUser validates, hashes and persists a user, then returns the stored record.
// The returned code is Success, Conflict or BadRequest for business outcomes.
func createUser(ctx context.Context, svcCtx *svc.ServiceContext, p createUserParams) (*user.UserInfo, int32, error) {
	username := strings.TrimSpace(p.Username)
	nickname := strings.TrimSpace(p.Nickname)
	email := strings.TrimSpace(p.Email)
	phone := strings.TrimSpace(p.Phone)
	avatar := strings.TrimSpace(p.Avatar)
	remark := strings.TrimSpace(p.Remark)

	if !validCreateUserInput(username, p.Password, nickname, email, phone, avatar, remark, p.Status, p.Role) {
		return nil, errcode.BadRequest, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errcode.ServerError, err
	}

	result, err := svcCtx.UserModel.Insert(ctx, &model.User{
		Username: username,
		Password: string(hash),
		Nickname: nickname,
		Email:    email,
		Phone:    phone,
		Avatar:   avatar,
		Status:   int64(p.Status),
		Remark:   remark,
		Role:     int64(p.Role),
	})
	if err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateKey {
			return nil, errcode.Conflict, nil
		}
		return nil, errcode.ServerError, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errcode.ServerError, err
	}

	stored, err := svcCtx.UserModel.FindOne(ctx, id)
	if err != nil {
		return nil, errcode.ServerError, err
	}
	return toUserInfo(stored), errcode.Success, nil
}

func validCreateUserInput(username, password, nickname, email, phone, avatar, remark string, status, role int32) bool {
	if utf8.RuneCountInString(username) < minUsernameLength || utf8.RuneCountInString(username) > maxUsernameLength {
		return false
	}
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return false
	}
	if utf8.RuneCountInString(nickname) > maxNicknameLength ||
		utf8.RuneCountInString(email) > maxEmailLength ||
		utf8.RuneCountInString(phone) > maxPhoneLength ||
		utf8.RuneCountInString(avatar) > maxAvatarLength ||
		utf8.RuneCountInString(remark) > maxRemarkLength {
		return false
	}
	if email != "" && !emailPattern.MatchString(email) {
		return false
	}
	return (status == 0 || status == 1) && (role == 0 || role == 1 || role == 2)
}
