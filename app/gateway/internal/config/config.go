package config

import "github.com/zeromicro/go-zero/gateway"

// JwtConfig holds token verification settings.
type JwtConfig struct {
	Secret string
}

type Config struct {
	gateway.GatewayConf
	Jwt JwtConfig
}
