package logic

import (
	"strings"
	"testing"
)

func TestGenerateNoIsUniqueAndPrefixed(t *testing.T) {
	first := generateNo("O")
	second := generateNo("O")
	if !strings.HasPrefix(first, "O") || !strings.HasPrefix(second, "O") {
		t.Fatalf("missing prefix: %q %q", first, second)
	}
	if first == second {
		t.Fatalf("expected unique order numbers, got %q twice", first)
	}
	if len(first) > 32 {
		t.Fatalf("order number too long: %d", len(first))
	}
}

func TestTrackNoValid(t *testing.T) {
	if !trackNoValid("顺丰", "SF123") {
		t.Fatal("expected valid tracking number")
	}
	if trackNoValid("顺丰", "") {
		t.Fatal("expected empty tracking number to be invalid")
	}
	if trackNoValid(strings.Repeat("a", maxLogisticsCompany+1), "SF123") {
		t.Fatal("expected oversized company to be invalid")
	}
	if trackNoValid("顺丰", strings.Repeat("a", maxTrackingNo+1)) {
		t.Fatal("expected oversized tracking number to be invalid")
	}
}
