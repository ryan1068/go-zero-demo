package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMiddlewareRejectsMissingToken(t *testing.T) {
	called := false
	handler := Middleware(testSecret)(func(http.ResponseWriter, *http.Request) {
		called = true
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodGet, "/order/1", nil))

	if called {
		t.Fatal("next handler should not be called")
	}
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}

	var body map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["code"] != float64(http.StatusUnauthorized) {
		t.Errorf("body code = %v, want %d", body["code"], http.StatusUnauthorized)
	}
}

func TestMiddlewareAcceptsBearerToken(t *testing.T) {
	token, err := Sign(testSecret, 7, "alice", time.Hour)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	var gotID int64
	var gotName string
	handler := Middleware(testSecret)(func(_ http.ResponseWriter, r *http.Request) {
		gotID, _ = UserID(r.Context())
		gotName, _ = Username(r.Context())
	})

	request := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	handler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if gotID != 7 || gotName != "alice" {
		t.Fatalf("identity = (%d, %q), want (7, alice)", gotID, gotName)
	}
}

func TestMiddlewareAcceptsRawTokenAndSetsMetadataHeaders(t *testing.T) {
	token, err := Sign(testSecret, 8, "bob", time.Hour)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	var userIDHeader string
	handler := Middleware(testSecret)(func(_ http.ResponseWriter, r *http.Request) {
		userIDHeader = r.Header.Get(metadataHeaderUserID)
	})

	request := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	request.Header.Set("Authorization", token)
	handler(httptest.NewRecorder(), request)

	if userIDHeader != "8" {
		t.Fatalf("metadata header = %q, want 8", userIDHeader)
	}
}

func TestMiddlewareSkipsLoginAndRemovesSpoofedHeaders(t *testing.T) {
	called := false
	handler := Middleware(testSecret, "/user/login")(func(_ http.ResponseWriter, r *http.Request) {
		called = true
		if got := r.Header.Get(metadataHeaderUserID); got != "" {
			t.Errorf("spoofed user id header was not removed: %q", got)
		}
	})

	request := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(`{}`))
	request.Header.Set(metadataHeaderUserID, "999")
	recorder := httptest.NewRecorder()
	handler(recorder, request)

	if !called {
		t.Fatal("next handler should be called for skipped path")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
}

func TestMiddlewareRejectsInvalidToken(t *testing.T) {
	called := false
	handler := Middleware(testSecret)(func(http.ResponseWriter, *http.Request) {
		called = true
	})

	request := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	request.Header.Set("Authorization", "Bearer not-a-token")
	recorder := httptest.NewRecorder()
	handler(recorder, request)

	if called {
		t.Fatal("next handler should not be called")
	}
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}
