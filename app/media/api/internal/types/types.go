package types

type UploadResp struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Url  string `json:"url"`
}

type UploadBatchResp struct {
	Code int32    `json:"code"`
	Msg  string   `json:"msg"`
	Urls []string `json:"urls"`
}

type BaseResp struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

type UploadNameReq struct {
	Name string `path:"name"`
}
