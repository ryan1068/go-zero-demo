package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ServeUploadHandler keeps the old /uploads/:name URLs working by proxying
// them to the media service. New uploads should use /api/media/files/:name.
func ServeUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	target, err := url.Parse(svcCtx.Config.Media.Url)
	if err != nil || target.Scheme == "" || target.Host == "" {
		logx.Errorf("invalid media service url: %q", svcCtx.Config.Media.Url)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "media service unavailable", http.StatusBadGateway)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
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

		r.URL.Path = strings.TrimRight(target.Path, "/") + "/files/" + name
		r.URL.RawQuery = ""
		proxy.ServeHTTP(w, r)
	}
}
