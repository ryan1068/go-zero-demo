package svc

import (
	"go-zero-demo/product-api/internal/config"
	"go-zero-demo/product-api/internal/middleware"
	"go-zero-demo/product-rpc/productclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	ProductRpc productclient.Product

	AdminOnlyMiddleware      rest.Middleware
	SuperAdminOnlyMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpcConf)),

		AdminOnlyMiddleware:      middleware.NewAdminOnlyMiddleware().Handle,
		SuperAdminOnlyMiddleware: middleware.NewSuperAdminOnlyMiddleware().Handle,
	}
}
