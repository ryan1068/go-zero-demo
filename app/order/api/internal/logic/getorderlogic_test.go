package logic

import (
	"context"
	"errors"
	"testing"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"google.golang.org/grpc"
)

type getOrderMock struct {
	createOrderFunc func(ctx context.Context, in *order.CreateOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error)
	getOrderFunc    func(ctx context.Context, in *order.GetOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error)
}

func (m *getOrderMock) CreateOrder(ctx context.Context, in *order.CreateOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
	return m.createOrderFunc(ctx, in, opts...)
}

func (m *getOrderMock) GetOrder(ctx context.Context, in *order.GetOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
	return m.getOrderFunc(ctx, in, opts...)
}

func TestGetOrder_Success(t *testing.T) {
	mock := &getOrderMock{
		getOrderFunc: func(ctx context.Context, in *order.GetOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
			return &order.CommonResp{
				Code: 200,
				Msg:  "ok",
				Data: "order_detail",
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{OrderRpc: mock}
	logic := NewGetOrderLogic(context.Background(), svcCtx)

	resp, err := logic.GetOrder(&types.GetOrderReq{Id: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if resp.Data != "order_detail" {
		t.Errorf("expected data 'order_detail', got '%s'", resp.Data)
	}
}

func TestGetOrder_RpcError(t *testing.T) {
	mock := &getOrderMock{
		getOrderFunc: func(ctx context.Context, in *order.GetOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
			return nil, errors.New("order not found")
		},
	}

	svcCtx := &svc.ServiceContext{OrderRpc: mock}
	logic := NewGetOrderLogic(context.Background(), svcCtx)

	_, err := logic.GetOrder(&types.GetOrderReq{Id: 999})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetOrder_EmptyId(t *testing.T) {
	mock := &getOrderMock{
		getOrderFunc: func(ctx context.Context, in *order.GetOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
			return &order.CommonResp{
				Code: 400,
				Msg:  "id required",
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{OrderRpc: mock}
	logic := NewGetOrderLogic(context.Background(), svcCtx)

	resp, err := logic.GetOrder(&types.GetOrderReq{Id: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 400 {
		t.Errorf("expected code 400, got %d", resp.Code)
	}
}
