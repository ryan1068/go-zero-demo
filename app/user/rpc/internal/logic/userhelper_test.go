package logic

import (
	"context"
	"strings"
	"testing"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

func TestValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"letters and digits", "abc12345", true},
		{"max length", strings.Repeat("a1", 36), true},
		{"chinese with letter and digit", "密码密码密码密码a1", true},
		{"too short", "ab12345", false},
		{"chinese only", "密码密码密码密码", false},
		{"letters only", "abcdefgh", false},
		{"digits only", "12345678", false},
		{"too long", strings.Repeat("a1", 37), false},
		{"too long in bytes", strings.Repeat("密", 24) + "a1", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validPassword(tt.password); got != tt.want {
				t.Fatalf("validPassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

func TestChangePassword_AllowsLegacyShortOldPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte("old123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash legacy password: %v", err)
	}
	rows := newUserRows()
	addUserRow(rows, 1, "alice", string(hash), "Alice", "", "", "", 1, "", 0)
	mock.ExpectQuery("select .* from `user` where `id` = \\? limit 1").
		WithArgs(int64(1)).
		WillReturnRows(rows)
	mock.ExpectExec("update `user` set `password`").
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(auth.MetadataUserIDKey, "1"))
	got, err := NewChangePasswordLogic(ctx, newUserServiceContext(t, db)).ChangePassword(&user.ChangePasswordReq{
		OldPassword: "old123",
		NewPassword: "newpass1",
	})
	if err != nil {
		t.Fatalf("ChangePassword() error = %v", err)
	}
	if got.Code != 200 {
		t.Fatalf("ChangePassword() code = %d, want 200", got.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestChangePassword_RejectsWeakNewPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"digits only", "12345678"},
		{"letters only", "abcdefgh"},
		{"too short", "abc1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(),
				metadata.Pairs(auth.MetadataUserIDKey, "1"))
			got, err := NewChangePasswordLogic(ctx, &svc.ServiceContext{}).ChangePassword(&user.ChangePasswordReq{
				OldPassword: "oldpass1",
				NewPassword: tt.password,
			})
			if err != nil {
				t.Fatalf("ChangePassword() error = %v", err)
			}
			if got.Code != 400 {
				t.Fatalf("ChangePassword() code = %d, want 400", got.Code)
			}
		})
	}
}
