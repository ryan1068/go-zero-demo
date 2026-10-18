package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/user-api/internal/logic"
	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
)

func ListRegionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListRegionsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		resp := logic.NewListRegionsLogic(r.Context(), svcCtx).ListRegions(&req)
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
