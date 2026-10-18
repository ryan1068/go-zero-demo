package main

import (
	"flag"
	"fmt"
	"net/http"

	"gateway/internal/config"
	"gateway/internal/middleware"
	"gateway/internal/types"

	"go-zero-demo/pkg/auth"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if err := auth.ValidateSecret(c.Jwt.Secret); err != nil {
		logx.Must(err)
	}

	gw := gateway.MustNewServer(c.GatewayConf)
	defer gw.Stop()

	// 注册全局中间件
	gw.Use(auth.Middleware(c.Jwt.Secret, "/user/login", "/product/"))
	gw.Use(middleware.NewLoggingMiddleware().Handle)

	// 注册自定义错误处理器
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *types.BizError:
			code := e.Code
			if code < http.StatusBadRequest || code > 599 {
				code = http.StatusBadRequest
			}
			return code, e
		default:
			logx.Errorf("gateway error: %v", err)
			return http.StatusInternalServerError, map[string]interface{}{
				"code": http.StatusInternalServerError,
				"msg":  "服务器内部错误",
			}
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	gw.Start()
}
