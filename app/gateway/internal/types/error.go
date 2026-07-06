package types

type BizError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *BizError) Error() string {
	return e.Msg
}

func NewBizError(code int, msg string) *BizError {
	return &BizError{
		Code: code,
		Msg:  msg,
	}
}