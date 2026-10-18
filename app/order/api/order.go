package main

import (
	"flag"
	"fmt"
	"net/http"

	"go-zero-demo/order-api/internal/config"
	"go-zero-demo/order-api/internal/handler"
	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/pkg/auth"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if err := auth.ValidateSecret(c.Jwt.Secret); err != nil {
		logx.Must(err)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(auth.Middleware(c.Jwt.Secret))
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return http.StatusUnauthorized, map[string]interface{}{"code": http.StatusUnauthorized, "msg": "未授权"}
			case codes.PermissionDenied:
				return http.StatusForbidden, map[string]interface{}{"code": http.StatusForbidden, "msg": "无权限"}
			case codes.NotFound:
				return http.StatusNotFound, map[string]interface{}{"code": http.StatusNotFound, "msg": "资源不存在"}
			case codes.InvalidArgument:
				return http.StatusBadRequest, map[string]interface{}{"code": http.StatusBadRequest, "msg": "请求参数错误"}
			default:
				logx.Errorf("order-api rpc error: %v", err)
				return http.StatusInternalServerError, map[string]interface{}{"code": http.StatusInternalServerError, "msg": "服务器内部错误"}
			}
		}

		logx.Errorf("order-api request error: %v", err)
		return http.StatusBadRequest, map[string]interface{}{"code": http.StatusBadRequest, "msg": "请求参数错误"}
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
