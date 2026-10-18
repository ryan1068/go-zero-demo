package logic

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	countUsersSQL = "(?i)select count\\(\\*\\) from .user. where"
	listUsersSQL  = "(?i)select .* from .user. where .* order by .id. desc limit \\? offset \\?"
)

func newListUsersLogic(t *testing.T, db *sql.DB) *ListUsersLogic {
	t.Helper()
	logic := NewListUsersLogic(context.Background(), newUserServiceContext(t, db))
	logic.Logger = logx.WithContext(context.Background())
	return logic
}

func TestListUsersRpc_Defaults(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(countUsersSQL).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))

	rows := addUserRow(newUserRows(), 100, "alice", "", "Alice", "alice@example.com", "13900000000", "", 1, "测试用户", 0)
	mock.ExpectQuery(listUsersSQL).
		WithArgs(int64(10), int64(0)).
		WillReturnRows(rows)

	resp, err := newListUsersLogic(t, db).ListUsers(&user.ListUsersReq{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Success {
		t.Fatalf("expected code %d, got %d", errcode.Success, resp.Code)
	}
	if resp.Page != 1 || resp.PageSize != defaultUserPageSize {
		t.Fatalf("unexpected pagination: %d/%d", resp.Page, resp.PageSize)
	}
	if resp.Total != 1 || len(resp.List) != 1 || resp.List[0].Username != "alice" {
		t.Fatalf("unexpected list: %+v", resp.List)
	}
}

func TestListUsersRpc_WithFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	args := []driver.Value{"adm%", "员%", "admin@%", "138%", int64(1)}
	mock.ExpectQuery(countUsersSQL).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	listArgs := append(append([]driver.Value{}, args...), int64(20), int64(20))
	mock.ExpectQuery(listUsersSQL).
		WithArgs(listArgs...).
		WillReturnRows(newUserRows())

	resp, err := newListUsersLogic(t, db).ListUsers(&user.ListUsersReq{
		Page:     2,
		PageSize: 20,
		Username: " adm ",
		Nickname: "员",
		Email:    "admin@",
		Phone:    "138",
		Status:   1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.Success || resp.Total != 0 || len(resp.List) != 0 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestListUsersRpc_FilterTooLong(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	resp, err := newListUsersLogic(t, db).ListUsers(&user.ListUsersReq{
		Username: strings.Repeat("a", maxUsernameLength+1),
	})
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

func TestListUsersRpc_CountError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(countUsersSQL).
		WillReturnError(errors.New("db down"))

	if _, err := newListUsersLogic(t, db).ListUsers(&user.ListUsersReq{}); err == nil {
		t.Fatal("expected database error, got nil")
	}
}
