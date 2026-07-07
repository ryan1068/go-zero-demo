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

type userLogicMock struct {
	userFunc func(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error)
}

func (m *userLogicMock) Login(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error) {
	return nil, nil
}

func (m *userLogicMock) User(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error) {
	return m.userFunc(ctx, in, opts...)
}

func TestUser_Success(t *testing.T) {
	mock := &userLogicMock{
		userFunc: func(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error) {
			return &user.UserResp{
				Message: "hello admin",
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{UserRpc: mock}
	logic := NewUserLogic(context.Background(), svcCtx)

	resp, err := logic.User(&types.UserReq{Name: "admin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Message != "hello admin" {
		t.Errorf("expected 'hello admin', got '%s'", resp.Message)
	}
}

func TestUser_RpcError(t *testing.T) {
	mock := &userLogicMock{
		userFunc: func(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error) {
			return nil, errors.New("rpc timeout")
		},
	}

	svcCtx := &svc.ServiceContext{UserRpc: mock}
	logic := NewUserLogic(context.Background(), svcCtx)

	_, err := logic.User(&types.UserReq{Name: "admin"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUser_EmptyName(t *testing.T) {
	mock := &userLogicMock{
		userFunc: func(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error) {
			return &user.UserResp{Message: ""}, nil
		},
	}

	svcCtx := &svc.ServiceContext{UserRpc: mock}
	logic := NewUserLogic(context.Background(), svcCtx)

	resp, err := logic.User(&types.UserReq{Name: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Message != "" {
		t.Errorf("expected empty message, got '%s'", resp.Message)
	}
}
