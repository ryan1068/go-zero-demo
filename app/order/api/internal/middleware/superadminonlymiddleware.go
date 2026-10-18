package middleware

import (
	"net/http"

	"go-zero-demo/pkg/auth"
)

// SuperAdminOnlyMiddleware allows only super admins.
type SuperAdminOnlyMiddleware struct{}

func NewSuperAdminOnlyMiddleware() *SuperAdminOnlyMiddleware {
	return &SuperAdminOnlyMiddleware{}
}

func (m *SuperAdminOnlyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := auth.Role(r.Context())
		if !ok || !auth.IsSuperAdmin(role) {
			writeForbidden(w, "需要超级管理员权限")
			return
		}
		next(w, r)
	}
}
