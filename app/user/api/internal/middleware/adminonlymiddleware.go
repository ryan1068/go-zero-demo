package middleware

import (
	"encoding/json"
	"net/http"

	"go-zero-demo/pkg/auth"
)

// AdminOnlyMiddleware allows super admins and operators.
type AdminOnlyMiddleware struct{}

func NewAdminOnlyMiddleware() *AdminOnlyMiddleware {
	return &AdminOnlyMiddleware{}
}

func (m *AdminOnlyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := auth.Role(r.Context())
		if !ok || !auth.IsAdmin(role) {
			writeForbidden(w, "无权限")
			return
		}
		next(w, r)
	}
}

func writeForbidden(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(map[string]any{"code": http.StatusForbidden, "msg": msg})
}
