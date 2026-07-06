package middleware

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoggingMiddleware struct {
}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

func (m *LoggingMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建响应记录器，以便获取响应状态码
		lrw := &loggingResponseWriter{w, http.StatusOK}

		// 记录请求信息
		logx.Infof("REQUEST [%s] %s %s",
			r.Header.Get("X-Request-Id"),
			r.Method,
			r.URL.Path)

		// 处理请求
		next(lrw, r)

		// 记录响应信息
		logx.Infof("RESPONSE [%s] %s %s %d %v",
			r.Header.Get("X-Request-Id"),
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			time.Since(start))
	}
}

// 自定义响应写入器，用于捕获状态码
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}