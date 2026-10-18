package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-demo/user-api/internal/config"
	"go-zero-demo/user-rpc/userclient"
)

type ServiceContext struct {
	Config config.Config

	UserRpc userclient.User
	Redis   *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
		Redis:   redis.MustNewRedis(c.Redis),
	}
}
