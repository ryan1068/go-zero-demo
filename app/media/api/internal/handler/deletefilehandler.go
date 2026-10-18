package handler

import (
	"net/http"

	"go-zero-demo/media-api/internal/logic"
	"go-zero-demo/media-api/internal/svc"
	"go-zero-demo/media-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadNameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		resp, err := logic.NewDeleteFileLogic(r.Context(), svcCtx.Storage, svcCtx.MediaFiles).Delete(req.Name)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
