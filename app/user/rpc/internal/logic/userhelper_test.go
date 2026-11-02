package logic

import (
	"context"
	"strings"
	"testing"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

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
		{"too short", "ab12345", false},
		{"letters only", "abcdefgh", false},
		{"digits only", "12345678", false},
		{"too long", strings.Repeat("a1", 37), false},
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
