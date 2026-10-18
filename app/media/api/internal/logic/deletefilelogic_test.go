package logic

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"go-zero-demo/media-api/internal/storage"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
)

func TestDeleteFile_InvalidName(t *testing.T) {
	dir := t.TempDir()
	if writeErr := os.WriteFile(filepath.Join(dir, "a.png"), pngBytes, 0o644); writeErr != nil {
		t.Fatalf("write file: %v", writeErr)
	}

	ctx := auth.WithUserRole(context.Background(), 1, "admin", auth.RoleAdmin)
	resp, err := NewDeleteFileLogic(ctx, storage.NewLocal(dir)).Delete("../a.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.BadRequest {
		t.Fatalf("expected code %d, got %d (%s)", errcode.BadRequest, resp.Code, resp.Msg)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "a.png")); statErr != nil {
		t.Fatalf("file should still exist: %v", statErr)
	}
}

func TestDeleteFile_Permissions(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		fileExists   bool
		wantCode     int32
		wantFileGone bool
	}{
		{
			name:         "admin deletes file",
			ctx:          auth.WithUserRole(context.Background(), 1, "admin", auth.RoleAdmin),
			fileExists:   true,
			wantCode:     errcode.Success,
			wantFileGone: true,
		},
		{
			name:         "normal user is forbidden",
			ctx:          auth.WithUserRole(context.Background(), 2, "buyer", 0),
			fileExists:   true,
			wantCode:     errcode.Forbidden,
			wantFileGone: false,
		},
		{
			name:         "missing file",
			ctx:          auth.WithUserRole(context.Background(), 1, "admin", auth.RoleAdmin),
			fileExists:   false,
			wantCode:     errcode.NotFound,
			wantFileGone: true,
		},
		{
			name:         "unauthenticated",
			ctx:          context.Background(),
			fileExists:   true,
			wantCode:     errcode.Unauthorized,
			wantFileGone: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.fileExists {
				if writeErr := os.WriteFile(filepath.Join(dir, "a.png"), pngBytes, 0o644); writeErr != nil {
					t.Fatalf("write file: %v", writeErr)
				}
			}

			resp, err := NewDeleteFileLogic(tc.ctx, storage.NewLocal(dir)).Delete("a.png")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Code != tc.wantCode {
				t.Fatalf("expected code %d, got %d (%s)", tc.wantCode, resp.Code, resp.Msg)
			}

			_, statErr := os.Stat(filepath.Join(dir, "a.png"))
			fileGone := os.IsNotExist(statErr)
			if fileGone != tc.wantFileGone {
				t.Fatalf("expected file gone=%v, got %v (err=%v)", tc.wantFileGone, fileGone, statErr)
			}
		})
	}
}
