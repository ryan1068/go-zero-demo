package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/pb/user"
)

const (
	userSelectColumns = "id, username, COALESCE(nickname, ''), COALESCE(email, ''), COALESCE(phone, ''), COALESCE(avatar, ''), status, COALESCE(remark, ''), UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at), role"
	mysqlDuplicateKey = 1062
)

type rowScanner interface {
	Scan(dest ...any) error
}

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

func scanUserInfo(row rowScanner) (*user.UserInfo, error) {
	info := new(user.UserInfo)
	err := row.Scan(
		&info.Id,
		&info.Username,
		&info.Nickname,
		&info.Email,
		&info.Phone,
		&info.Avatar,
		&info.Status,
		&info.Remark,
		&info.CreatedAt,
		&info.UpdatedAt,
		&info.Role,
	)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func findUserByID(ctx context.Context, db *sql.DB, id int64) (*user.UserInfo, error) {
	return scanUserInfo(db.QueryRowContext(ctx,
		"SELECT "+userSelectColumns+" FROM user WHERE id = ?", id,
	))
}

// createUser inserts a user and returns the persisted record. The returned
// code is Success, Conflict or BadRequest when the failure is a business error.
func createUser(ctx context.Context, db *sql.DB, p createUserParams) (*user.UserInfo, int32, error) {
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

	result, err := db.ExecContext(ctx,
		"INSERT INTO user (username, password, nickname, email, phone, avatar, status, role, remark) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		username, string(hash), nickname, email, phone, avatar, p.Status, p.Role, remark,
	)
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

	info, err := findUserByID(ctx, db, id)
	if err != nil {
		return nil, errcode.ServerError, err
	}
	return info, errcode.Success, nil
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
