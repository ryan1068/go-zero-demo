package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type JwtConfig struct {
	Secret string
}

type UploadConfig struct {
	Dir       string
	URLPrefix string
	MaxSizeMB int64
}

type MediaConfig struct {
	Url string
}

type Config struct {
	rest.RestConf
	ProductRpcConf zrpc.RpcClientConf
	Jwt            JwtConfig
	Upload         UploadConfig
	Media          MediaConfig
}
