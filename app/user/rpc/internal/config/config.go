package config

import "github.com/zeromicro/go-zero/zrpc"

// JwtConfig holds token signing settings.
type JwtConfig struct {
	Secret      string
	ExpireHours int
}

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Jwt JwtConfig
}
