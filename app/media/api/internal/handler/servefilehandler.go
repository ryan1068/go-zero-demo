package handler

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-zero-demo/media-api/internal/svc"
	"go-zero-demo/media-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ServeFileHandler streams an uploaded file. goctl cannot generate
// non-JSON responses, so this handler is intentionally hand-written.
func ServeFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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

		file, err := svcCtx.Storage.Open(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				http.Error(w, "close file failed", http.StatusInternalServerError)
			}
		}()

		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeContent(w, r, name, time.Time{}, file)
	}
}
