package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/product-api/internal/config"
	"go-zero-demo/product-api/internal/handler"
	"go-zero-demo/product-api/internal/svc"
)

var configFile = flag.String("f", "etc/product.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if err := auth.ValidateSecret(c.Jwt.Secret); err != nil {
		logx.Must(err)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9104", nil); err != nil {
			logx.Errorf("metrics server failed: %v", err)
		}
	}()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(auth.Middleware(c.Jwt.Secret, "/categories", "/products", "/products/", "/product/", "/uploads/"))
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		logx.Errorf("product-api request error: %v", err)
		return http.StatusInternalServerError, map[string]interface{}{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部错误",
		}
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
