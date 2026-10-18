package svc

import (
	"go-zero-demo/media-api/internal/config"
	"go-zero-demo/media-api/internal/storage"
)

type ServiceContext struct {
	Config  config.Config
	Storage storage.Storage
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Storage: storage.NewLocal(c.Upload.Dir),
	}
}
