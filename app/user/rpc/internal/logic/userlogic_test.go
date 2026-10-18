package logic

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zeromicro/go-zero/core/logx"
)

const findUserByUsernameSQL = "(?i)select .* from .user. where .username. = \\? limit 1"

func TestUserRpc_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	rows := addUserRow(newUserRows(), 1, "admin", "", "管理员", "admin@example.com", "13800000000", "", 1, "系统管理员", 1)
	mock.ExpectQuery(findUserByUsernameSQL).
		WithArgs("admin").
		WillReturnRows(rows)

	logic := NewUserLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.User(&user.UserReq{Name: "admin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Success {
		t.Fatalf("expected code %d, got %d", errcode.Success, resp.Code)
	}
	if resp.Data == nil {
		t.Fatal("expected user data, got nil")
	}
	if resp.Data.Id != 1 || resp.Data.Username != "admin" || resp.Data.Nickname != "管理员" {
		t.Fatalf("unexpected user: %+v", resp.Data)
	}
	if resp.Data.Email != "admin@example.com" || resp.Data.Phone != "13800000000" || resp.Data.Status != 1 || resp.Data.Role != 1 {
		t.Fatalf("unexpected extra fields: %+v", resp.Data)
	}
}

func TestUserRpc_EmptyName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	logic := NewUserLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.User(&user.UserReq{Name: "   "})
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

func TestUserRpc_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(findUserByUsernameSQL).
		WithArgs("ghost").
		WillReturnError(sql.ErrNoRows)

	logic := NewUserLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.User(&user.UserReq{Name: "ghost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.NotFound {
		t.Fatalf("expected code %d, got %d", errcode.NotFound, resp.Code)
	}
}

func TestUserRpc_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(findUserByUsernameSQL).
		WithArgs("admin").
		WillReturnError(errors.New("db down"))

	logic := NewUserLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())

	if _, err := logic.User(&user.UserReq{Name: "admin"}); err == nil {
		t.Fatal("expected database error, got nil")
	}
}
