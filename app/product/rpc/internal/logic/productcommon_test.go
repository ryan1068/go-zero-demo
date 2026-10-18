package logic

import (
	"strings"
	"testing"
)

func TestEncodeImages(t *testing.T) {
	got, err := encodeImages(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[]" {
		t.Fatalf("expected empty JSON array, got %q", got)
	}
	got, err = encodeImages([]string{"/a.png", "/b.png"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != `["/a.png","/b.png"]` {
		t.Fatalf("unexpected JSON: %q", got)
	}
}

func TestValidText(t *testing.T) {
	if !validText("商品", 4) {
		t.Fatal("expected text within limit")
	}
	if validText(strings.Repeat("a", 5), 4) {
		t.Fatal("expected text over limit to be invalid")
	}
}
