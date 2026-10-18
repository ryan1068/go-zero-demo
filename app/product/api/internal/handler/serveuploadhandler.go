package handler

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
)

// ServeUploadHandler streams an uploaded file. goctl cannot generate
// non-JSON responses, so this handler is intentionally hand-written.
func ServeUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadNameReq
		if err := httpx.Parse(r, &req); err != nil {
			http.NotFound(w, r)
			return
		}
		name := filepath.Base(req.Name)
		if name != req.Name || strings.Contains(name, "..") || name == "." {
			http.NotFound(w, r)
			return
		}
		dir := svcCtx.Config.Upload.Dir
		if dir == "" {
			dir = "./data/uploads"
		}
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, filepath.Join(dir, name))
	}
}
