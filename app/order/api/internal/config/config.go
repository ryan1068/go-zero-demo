package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// JwtConfig holds token verification settings.
type JwtConfig struct {
	Secret string
}

type Config struct {
	rest.RestConf

	OrderRpcConf zrpc.RpcClientConf
	Jwt          JwtConfig
}
