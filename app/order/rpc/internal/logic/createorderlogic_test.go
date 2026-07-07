package logic

import (
	"context"
	"errors"
	"testing"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/user-rpc/pb/user"

	"google.golang.org/grpc"
)

type rpcUserMock struct {
	userFunc func(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error)
}

func (m *rpcUserMock) Login(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error) {
	return nil, nil
}

func (m *rpcUserMock) User(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error) {
	return m.userFunc(ctx, in, opts...)
}

func TestCreateOrderRpc_Success(t *testing.T) {
	mock := &rpcUserMock{
		userFunc: func(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error) {
			return &user.UserResp{Message: "goods"}, nil
		},
	}

	svcCtx := &svc.ServiceContext{UserRpc: mock}
	logic := NewCreateOrderLogic(context.Background(), svcCtx)

	resp, err := logic.CreateOrder(&order.CreateOrderReq{
		GoodsId:   1,
		GoodsName: "book",
		Price:     9990,
		CreatedAt: 1704067200,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if resp.Data != "book" {
		t.Errorf("expected data 'book', got '%s'", resp.Data)
	}
}

func TestCreateOrderRpc_UserRpcError(t *testing.T) {
	mock := &rpcUserMock{
		userFunc: func(ctx context.Context, in *user.UserReq, opts ...grpc.CallOption) (*user.UserResp, error) {
			return nil, errors.New("user service unavailable")
		},
	}

	svcCtx := &svc.ServiceContext{UserRpc: mock}
	logic := NewCreateOrderLogic(context.Background(), svcCtx)

	resp, err := logic.CreateOrder(&order.CreateOrderReq{
		GoodsId:   1,
		GoodsName: "book",
		Price:     9990,
		CreatedAt: 1704067200,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200 even on user rpc error, got %d", resp.Code)
	}
}
