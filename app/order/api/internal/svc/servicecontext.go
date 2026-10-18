package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-demo/order-api/internal/config"
	"go-zero-demo/order-api/internal/middleware"
	"go-zero-demo/order-rpc/orderclient"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc orderclient.Order

	AdminOnlyMiddleware      rest.Middleware
	SuperAdminOnlyMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),

		AdminOnlyMiddleware:      middleware.NewAdminOnlyMiddleware().Handle,
		SuperAdminOnlyMiddleware: middleware.NewSuperAdminOnlyMiddleware().Handle,
	}
}
