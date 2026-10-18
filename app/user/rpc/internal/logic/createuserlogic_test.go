package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	insertUserSQL   = "(?i)insert into .user."
	findUserByIDSQL = "(?i)select .* from .user. where .id. = \\? limit 1"
	insertAuditSQL  = "(?i)insert into .admin_audit_log."
)

func validCreateUserReq() *user.CreateUserReq {
	return &user.CreateUserReq{
		Username: "alice",
		Password: "secret123",
		Nickname: "Alice",
		Email:    "alice@example.com",
		Phone:    "13900000000",
		Avatar:   "/avatars/alice.svg",
		Status:   1,
		Role:     0,
		Remark:   "测试用户",
	}
}

func TestCreateUserRpc_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(insertUserSQL).
		WithArgs("alice", sqlmock.AnyArg(), "Alice", "alice@example.com", "13900000000", "/avatars/alice.svg", int64(1), "测试用户", int64(0), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(100, 1))

	rows := addUserRow(newUserRows(), 100, "alice", "", "Alice", "alice@example.com", "13900000000", "/avatars/alice.svg", 1, "测试用户", 0)
	mock.ExpectQuery(findUserByIDSQL).
		WithArgs(int64(100)).
		WillReturnRows(rows)
	mock.ExpectExec(insertAuditSQL).
		WithArgs(int64(0), "create_user", "user", int64(100), "alice").
		WillReturnResult(sqlmock.NewResult(1, 1))

	logic := NewCreateUserLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.CreateUser(validCreateUserReq())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Success {
		t.Fatalf("expected code %d, got %d", errcode.Success, resp.Code)
	}
	if resp.Data == nil || resp.Data.Id != 100 || resp.Data.Username != "alice" {
		t.Fatalf("unexpected user: %+v", resp.Data)
	}
	if resp.Data.Email != "alice@example.com" || resp.Data.Status != 1 {
		t.Fatalf("unexpected extra fields: %+v", resp.Data)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCreateUserRpc_DuplicateUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(insertUserSQL).
		WillReturnError(&mysqlDriver.MySQLError{Number: mysqlDuplicateKey, Message: "Duplicate entry"})

	logic := NewCreateUserLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.CreateUser(validCreateUserReq())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Conflict {
		t.Fatalf("expected code %d, got %d", errcode.Conflict, resp.Code)
	}
}

func TestCreateUserRpc_Validation(t *testing.T) {
	tests := []struct {
		name string
		req  *user.CreateUserReq
	}{
		{"short username", &user.CreateUserReq{Username: "ab", Password: "secret123", Status: 1}},
		{"long username", &user.CreateUserReq{Username: strings.Repeat("a", maxUsernameLength+1), Password: "secret123", Status: 1}},
		{"short password", &user.CreateUserReq{Username: "alice", Password: "12345", Status: 1}},
		{"invalid email", &user.CreateUserReq{Username: "alice", Password: "secret123", Email: "not-an-email", Status: 1}},
		{"invalid status", &user.CreateUserReq{Username: "alice", Password: "secret123", Status: 2}},
		{"nickname too long", &user.CreateUserReq{Username: "alice", Password: "secret123", Nickname: strings.Repeat("a", maxNicknameLength+1), Status: 1}},
		{"avatar too long", &user.CreateUserReq{Username: "alice", Password: "secret123", Avatar: strings.Repeat("a", maxAvatarLength+1), Status: 1}},
		{"remark too long", &user.CreateUserReq{Username: "alice", Password: "secret123", Remark: strings.Repeat("a", maxRemarkLength+1), Status: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			logic := NewCreateUserLogic(context.Background(), newUserServiceContext(t, db))
			logic.Logger = logx.WithContext(context.Background())

			resp, err := logic.CreateUser(tt.req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Code != errcode.BadRequest {
				t.Fatalf("expected code %d, got %d", errcode.BadRequest, resp.Code)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unexpected query: %v", err)
			}
		})
	}
}

func TestCreateUserRpc_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(insertUserSQL).
		WillReturnError(errors.New("db down"))

	logic := NewCreateUserLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	if _, err := logic.CreateUser(validCreateUserReq()); err == nil {
		t.Fatal("expected database error, got nil")
	}
}

func TestCreateUserRpc_FindError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(insertUserSQL).
		WillReturnResult(sqlmock.NewResult(100, 1))
	mock.ExpectQuery(findUserByIDSQL).
		WithArgs(int64(100)).
		WillReturnError(sql.ErrNoRows)

	logic := NewCreateUserLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	if _, err := logic.CreateUser(validCreateUserReq()); err == nil {
		t.Fatal("expected query error, got nil")
	}
}
