package main

import (
	"flag"
	"fmt"
	"gateway/internal/middleware"
	"gateway/internal/types"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")

func main() {
	flag.Parse()

	var c gateway.GatewayConf
	conf.MustLoad(*configFile, &c)
	gw := gateway.MustNewServer(c)
	defer gw.Stop()

	// 注册全局中间件
	gw.Use(middleware.NewAuthMiddleware().Handle)
	gw.Use(middleware.NewLoggingMiddleware().Handle)

	// 注册自定义错误处理器
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *types.BizError:
			return http.StatusBadRequest, e
		default:
			return http.StatusInternalServerError, map[string]interface{}{
				"code": http.StatusInternalServerError,
				"msg":  e.Error(),
			}
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	gw.Start()
}