package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalStorage_SaveOpenDelete(t *testing.T) {
	dir := t.TempDir()
	store := NewLocal(dir)
	ctx := context.Background()

	written, err := store.Save(ctx, "a.txt", strings.NewReader("hello"), 10)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if written != 5 {
		t.Fatalf("expected 5 bytes, got %d", written)
	}

	content, err := os.ReadFile(filepath.Join(dir, "a.txt"))
	if err != nil {
		t.Fatalf("read saved file: %v", err)
	}
	if string(content) != "hello" {
		t.Fatalf("expected hello, got %q", content)
	}

	f, err := store.Open("a.txt")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	opened, err := io.ReadAll(f)
	if closeErr := f.Close(); closeErr != nil {
		t.Fatalf("close: %v", closeErr)
	}
	if err != nil {
		t.Fatalf("read opened file: %v", err)
	}
	if string(opened) != "hello" {
		t.Fatalf("expected hello, got %q", opened)
	}

	if err := store.Delete("a.txt"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "a.txt")); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, got %v", err)
	}
}

func TestLocalStorage_RejectsTooLarge(t *testing.T) {
	dir := t.TempDir()
	store := NewLocal(dir)

	if _, err := store.Save(context.Background(), "b.txt", strings.NewReader("123456"), 5); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("expected ErrTooLarge, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "b.txt")); !os.IsNotExist(err) {
		t.Fatalf("oversized file should not exist, got %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected no temp files, got %d", len(entries))
	}
}

func TestLocalStorage_OpenUsesBaseName(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "secret"), []byte("inside"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	f, err := NewLocal(dir).Open("../secret")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(content) != "inside" {
		t.Fatalf("expected inside, got %q", content)
	}
}
