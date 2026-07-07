package logic

import (
	"context"
	"testing"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginRpc_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(1, "admin", string(hash))
	mock.ExpectQuery("SELECT id, username, password FROM user WHERE username = \\?").
		WithArgs("admin").
		WillReturnRows(rows)

	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewLoginLogic(context.Background(), svcCtx)
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.Login(&user.LoginReq{
		Username: "admin",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Success {
		t.Errorf("expected code %d, got %d", errcode.Success, resp.Code)
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

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(1, "admin", string(hash))
	mock.ExpectQuery("SELECT id, username, password FROM user WHERE username = \\?").
		WithArgs("admin").
		WillReturnRows(rows)

	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewLoginLogic(context.Background(), svcCtx)
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.Login(&user.LoginReq{
		Username: "admin",
		Password: "wrong",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Unauthorized {
		t.Errorf("expected code %d, got %d", errcode.Unauthorized, resp.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestLoginRpc_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, username, password FROM user WHERE username = \\?").
		WithArgs("unknown").
		WillReturnError(sqlmock.ErrCancelled)

	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewLoginLogic(context.Background(), svcCtx)
	logic.Logger = logx.WithContext(context.Background())

	resp, err := logic.Login(&user.LoginReq{
		Username: "unknown",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Unauthorized {
		t.Errorf("expected code %d, got %d", errcode.Unauthorized, resp.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
