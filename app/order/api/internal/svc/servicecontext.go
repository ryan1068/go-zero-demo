package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-demo/order-api/internal/config"
	"go-zero-demo/order-rpc/orderclient"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
	}
}
