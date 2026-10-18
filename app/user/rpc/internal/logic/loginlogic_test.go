package logic

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

const updateLastLoginSQL = "(?i)update .user. set .last_login_at."

func TestLoginRpc_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	rows := addUserRow(newUserRows(), 1, "admin", string(hash), "", "", "", "", 1, "", 1)
	mock.ExpectQuery(findUserByUsernameSQL).
		WithArgs("admin").
		WillReturnRows(rows)
	mock.ExpectExec(updateLastLoginSQL).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	logic := NewLoginLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.Login(&user.LoginReq{
		Username: "admin",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Success {
		t.Fatalf("expected code %d, got %d", errcode.Success, resp.Code)
	}
	if resp.Token == "" {
		t.Fatal("expected token, got empty string")
	}
	if resp.UserId != 1 || resp.Username != "admin" || resp.Role != 1 {
		t.Fatalf("unexpected identity: (%d, %q)", resp.UserId, resp.Username)
	}

	claims, err := auth.Parse(testSecret, resp.Token)
	if err != nil {
		t.Fatalf("token is not parseable: %v", err)
	}
	if claims.UserID != 1 || claims.Username != "admin" || claims.Role != 1 {
		t.Fatalf("unexpected claims: (%d, %q)", claims.UserID, claims.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestLoginRpc_WrongPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	rows := addUserRow(newUserRows(), 1, "admin", string(hash), "", "", "", "", 1, "", 0)
	mock.ExpectQuery(findUserByUsernameSQL).
		WithArgs("admin").
		WillReturnRows(rows)

	logic := NewLoginLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.Login(&user.LoginReq{
		Username: "admin",
		Password: "wrong",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Unauthorized {
		t.Fatalf("expected code %d, got %d", errcode.Unauthorized, resp.Code)
	}
	if resp.Token != "" {
		t.Fatalf("expected empty token, got %q", resp.Token)
	}
}

func TestLoginRpc_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(findUserByUsernameSQL).
		WithArgs("unknown").
		WillReturnError(sql.ErrNoRows)

	logic := NewLoginLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.Login(&user.LoginReq{
		Username: "unknown",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Unauthorized {
		t.Fatalf("expected code %d, got %d", errcode.Unauthorized, resp.Code)
	}
}

func TestLoginRpc_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(findUserByUsernameSQL).
		WithArgs("admin").
		WillReturnError(errors.New("db down"))

	logic := NewLoginLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	if _, err := logic.Login(&user.LoginReq{Username: "admin", Password: "123456"}); err == nil {
		t.Fatal("expected database error, got nil")
	}
}

func TestLoginRpc_EmptyUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	logic := NewLoginLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.Login(&user.LoginReq{Username: "  ", Password: "123456"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.BadRequest {
		t.Fatalf("expected code %d, got %d", errcode.BadRequest, resp.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected query: %v", err)
	}
}
