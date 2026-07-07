package logic

import (
	"context"
	"testing"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
)

func TestGetOrderRpc_Success(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	logic := NewGetOrderLogic(context.Background(), svcCtx)

	resp, err := logic.GetOrder(&order.GetOrderReq{Id: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if resp.Data != "1111" {
		t.Errorf("expected data '1111', got '%s'", resp.Data)
	}
}

func TestGetOrderRpc_ZeroId(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	logic := NewGetOrderLogic(context.Background(), svcCtx)

	resp, err := logic.GetOrder(&order.GetOrderReq{Id: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
}
