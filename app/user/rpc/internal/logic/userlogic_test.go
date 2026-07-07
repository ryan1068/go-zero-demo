package logic

import (
	"context"
	"testing"

	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"
)

func TestUserRpc_Success(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	logic := NewUserLogic(context.Background(), svcCtx)

	resp, err := logic.User(&user.UserReq{Name: "admin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Message != "admin345" {
		t.Errorf("expected 'admin345', got '%s'", resp.Message)
	}
}

func TestUserRpc_EmptyName(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	logic := NewUserLogic(context.Background(), svcCtx)

	resp, err := logic.User(&user.UserReq{Name: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Message != "345" {
		t.Errorf("expected '345', got '%s'", resp.Message)
	}
}
