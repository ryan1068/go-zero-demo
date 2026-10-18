package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"go-zero-demo/pkg/auth"
)

func adminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := auth.Role(r.Context())
		if !ok || !auth.IsAdmin(role) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(map[string]any{"code": http.StatusForbidden, "msg": "无权限"}); err != nil {
				log.Printf("failed to write forbidden response: %v", err)
			}
			return
		}
		next(w, r)
	}
}

func superAdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := auth.Role(r.Context())
		if !ok || !auth.IsSuperAdmin(role) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(map[string]any{"code": http.StatusForbidden, "msg": "需要超级管理员权限"}); err != nil {
				log.Printf("failed to write forbidden response: %v", err)
			}
			return
		}
		next(w, r)
	}
}
