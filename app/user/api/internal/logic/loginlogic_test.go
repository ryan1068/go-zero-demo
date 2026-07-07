package logic

import (
	"context"
	"errors"
	"testing"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"google.golang.org/grpc"
)

type mockUserRpc struct {
	loginFunc func(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error)
	userFunc  func(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error)
}

func (m *mockUserRpc) Login(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error) {
	return m.loginFunc(ctx, in, opts...)
}

func (m *mockUserRpc) User(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error) {
	return m.userFunc(ctx, in, opts...)
}

func TestLogin_Success(t *testing.T) {
	mock := &mockUserRpc{
		loginFunc: func(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error) {
			return &user.LoginResp{
				Code: 200,
				Msg:  "ok",
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{UserRpc: mock}
	logic := NewLoginLogic(context.Background(), svcCtx)

	resp, err := logic.Login(&types.LoginReq{
		Username: "admin",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if resp.Token != "token_admin" {
		t.Errorf("expected token 'token_admin', got '%s'", resp.Token)
	}
}

func TestLogin_RpcError(t *testing.T) {
	mock := &mockUserRpc{
		loginFunc: func(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error) {
			return nil, errors.New("rpc connection failed")
		},
	}

	svcCtx := &svc.ServiceContext{UserRpc: mock}
	logic := NewLoginLogic(context.Background(), svcCtx)

	_, err := logic.Login(&types.LoginReq{
		Username: "admin",
		Password: "123456",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mock := &mockUserRpc{
		loginFunc: func(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error) {
			return &user.LoginResp{
				Code: 401,
				Msg:  "invalid credentials",
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{UserRpc: mock}
	logic := NewLoginLogic(context.Background(), svcCtx)

	resp, err := logic.Login(&types.LoginReq{
		Username: "admin",
		Password: "wrong",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 401 {
		t.Errorf("expected code 401, got %d", resp.Code)
	}
	if resp.Token != "" {
		t.Errorf("expected empty token for failed login, got '%s'", resp.Token)
	}
}