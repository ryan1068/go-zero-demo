package errcode

const (
	Success      = 200
	BadRequest   = 400
	Unauthorized = 401
	Forbidden    = 403
	NotFound     = 404
	ServerError  = 500
)

var msgMap = map[int32]string{
	Success:      "ok",
	BadRequest:   "请求参数错误",
	Unauthorized: "未授权",
	Forbidden:    "无权限",
	NotFound:     "资源不存在",
	ServerError:  "服务器内部错误",
}

func Msg(code int32) string {
	if msg, ok := msgMap[code]; ok {
		return msg
	}
	return "未知错误"
}
