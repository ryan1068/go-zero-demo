package logic

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-zero-demo/media-api/internal/config"
	"go-zero-demo/media-api/internal/storage"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
)

var pngBytes = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")

func newUploadFile(t *testing.T, name string, content []byte) (multipart.File, *multipart.FileHeader) {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		t.Fatalf("parse multipart form: %v", err)
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		t.Fatalf("get form file: %v", err)
	}
	return file, header
}

func newUploadLogic(ctx context.Context, dir string) *UploadLogic {
	cfg := config.UploadConfig{Dir: dir, URLPrefix: "/api/media/files/", MaxSizeMB: 1}
	return NewUploadLogic(ctx, cfg, storage.NewLocal(dir))
}

func TestUpload_Validation(t *testing.T) {
	adminCtx := auth.WithUserRole(context.Background(), 7, "admin", auth.RoleAdmin)

	cases := []struct {
		name     string
		filename string
		content  []byte
		ctx      context.Context
		wantCode int32
	}{
		{name: "bad extension", filename: "a.txt", content: pngBytes, ctx: adminCtx, wantCode: errcode.BadRequest},
		{name: "too large", filename: "a.png", content: bytes.Repeat([]byte{0}, 1024*1024+1), ctx: adminCtx, wantCode: errcode.BadRequest},
		{name: "not an image", filename: "a.png", content: []byte("hello"), ctx: adminCtx, wantCode: errcode.BadRequest},
		{name: "empty file", filename: "a.png", content: nil, ctx: adminCtx, wantCode: errcode.BadRequest},
		{name: "unauthenticated", filename: "a.png", content: pngBytes, ctx: context.Background(), wantCode: errcode.Unauthorized},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			file, header := newUploadFile(t, tc.filename, tc.content)
			defer file.Close()

			resp, err := newUploadLogic(tc.ctx, dir).Save(file, header)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Code != tc.wantCode {
				t.Fatalf("expected code %d, got %d (%s)", tc.wantCode, resp.Code, resp.Msg)
			}

			entries, readErr := os.ReadDir(dir)
			if readErr != nil {
				t.Fatalf("read upload dir: %v", readErr)
			}
			if len(entries) != 0 {
				t.Fatalf("expected no files on disk, got %d", len(entries))
			}
		})
	}
}

func TestUploadBatch_Success(t *testing.T) {
	dir := t.TempDir()
	ctx := auth.WithUserRole(context.Background(), 7, "admin", auth.RoleAdmin)
	logic := newUploadLogic(ctx, dir)

	first, firstHeader := newUploadFile(t, "a.png", pngBytes)
	defer first.Close()
	second, secondHeader := newUploadFile(t, "b.png", pngBytes)
	defer second.Close()

	resp, err := logic.SaveBatch([]*multipart.FileHeader{firstHeader, secondHeader})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Fatalf("expected code 200, got %d (%s)", resp.Code, resp.Msg)
	}
	if len(resp.Urls) != 2 {
		t.Fatalf("expected 2 urls, got %d", len(resp.Urls))
	}
	for _, url := range resp.Urls {
		if !strings.HasPrefix(url, "/api/media/files/") {
			t.Fatalf("expected url prefix %q, got %q", "/api/media/files/", url)
		}
		name := strings.TrimPrefix(url, "/api/media/files/")
		if _, statErr := os.Stat(filepath.Join(dir, name)); statErr != nil {
			t.Fatalf("uploaded file not found: %v", statErr)
		}
	}
}

func TestUpload_Success(t *testing.T) {
	dir := t.TempDir()
	ctx := auth.WithUserRole(context.Background(), 7, "admin", auth.RoleAdmin)
	logic := newUploadLogic(ctx, dir)

	file, header := newUploadFile(t, "a.png", pngBytes)
	defer file.Close()

	resp, err := logic.Save(file, header)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != 200 {
		t.Fatalf("expected code 200, got %d (%s)", resp.Code, resp.Msg)
	}
	if !strings.HasPrefix(resp.Url, "/api/media/files/") {
		t.Fatalf("expected url prefix %q, got %q", "/api/media/files/", resp.Url)
	}

	name := strings.TrimPrefix(resp.Url, "/api/media/files/")
	if name == "" || strings.ContainsAny(name, "/\\") {
		t.Fatalf("unexpected file name %q", name)
	}
	if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
		t.Fatalf("uploaded file not found: %v", err)
	}
}
