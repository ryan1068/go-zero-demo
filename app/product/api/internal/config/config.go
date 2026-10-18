package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type JwtConfig struct {
	Secret string
}

type MediaConfig struct {
	Url string
}

type Config struct {
	rest.RestConf
	ProductRpcConf zrpc.RpcClientConf
	Jwt            JwtConfig
	Media          MediaConfig
}
