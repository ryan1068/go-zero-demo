package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type JwtConfig struct {
	Secret string
}

type UploadConfig struct {
	Dir       string `json:",default=./data/uploads"`
	URLPrefix string `json:",default=/api/media/files/"`
	MaxSizeMB int64  `json:",default=10"`
}

type Config struct {
	rest.RestConf
	Jwt    JwtConfig
	Upload UploadConfig
}
