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

type orderRpcMock struct {
	createOrderFunc func(ctx context.Context, in *order.CreateOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error)
	getOrderFunc    func(ctx context.Context, in *order.GetOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error)
}

func (m *orderRpcMock) CreateOrder(ctx context.Context, in *order.CreateOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
	return m.createOrderFunc(ctx, in, opts...)
}

func (m *orderRpcMock) GetOrder(ctx context.Context, in *order.GetOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
	return m.getOrderFunc(ctx, in, opts...)
}

func TestCreateOrder_Success(t *testing.T) {
	mock := &orderRpcMock{
		createOrderFunc: func(ctx context.Context, in *order.CreateOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
			return &order.CommonResp{
				Code: 200,
				Msg:  "ok",
				Data: "order_id_001",
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{OrderRpc: mock}
	logic := NewCreateOrderLogic(context.Background(), svcCtx)

	resp, err := logic.CreateOrder(&types.CreateOrderReq{
		Goods_id:   1,
		Goods_name: "book",
		Price:      9990,
		Created_at: 1704067200,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if resp.Data != "order_id_001" {
		t.Errorf("expected data 'order_id_001', got '%s'", resp.Data)
	}
}

func TestCreateOrder_RpcError(t *testing.T) {
	mock := &orderRpcMock{
		createOrderFunc: func(ctx context.Context, in *order.CreateOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
			return nil, errors.New("rpc unavailable")
		},
	}

	svcCtx := &svc.ServiceContext{OrderRpc: mock}
	logic := NewCreateOrderLogic(context.Background(), svcCtx)

	_, err := logic.CreateOrder(&types.CreateOrderReq{
		Goods_id:   1,
		Goods_name: "book",
		Price:      9990,
		Created_at: 1704067200,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateOrder_ZeroPrice(t *testing.T) {
	mock := &orderRpcMock{
		createOrderFunc: func(ctx context.Context, in *order.CreateOrderReq, opts ...grpc.CallOption) (*order.CommonResp, error) {
			return &order.CommonResp{
				Code: 400,
				Msg:  "invalid price",
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{OrderRpc: mock}
	logic := NewCreateOrderLogic(context.Background(), svcCtx)

	resp, err := logic.CreateOrder(&types.CreateOrderReq{
		Goods_id:   1,
		Goods_name: "book",
		Price:      0,
		Created_at: 1704067200,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 400 {
		t.Errorf("expected code 400, got %d", resp.Code)
	}
}
