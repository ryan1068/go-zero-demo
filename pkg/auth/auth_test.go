package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func TestValidateSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		want   error
	}{
		{"valid", testSecret, nil},
		{"empty", "", ErrInvalidSecret},
		{"too short", "short", ErrInvalidSecret},
		{"whitespace", "                                ", ErrInvalidSecret},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSecret(tt.secret); !errors.Is(err, tt.want) {
				t.Fatalf("ValidateSecret() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSignAndParse(t *testing.T) {
	token, err := Sign(testSecret, 42, "admin", time.Hour)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	if token == "" {
		t.Fatal("Sign() returned empty token")
	}

	claims, err := Parse(testSecret, token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("Username = %q, want admin", claims.Username)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt is nil")
	}
}

func TestSignRejectsInvalidInputs(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		userID int64
		ttl    time.Duration
		want   error
	}{
		{"short secret", "short", 1, time.Hour, ErrInvalidSecret},
		{"zero user id", testSecret, 0, time.Hour, ErrInvalidUserID},
		{"zero ttl", testSecret, 1, 0, ErrInvalidTTL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Sign(tt.secret, tt.userID, "admin", tt.ttl)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Sign() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestParseExpiredToken(t *testing.T) {
	token, err := Sign(testSecret, 1, "admin", time.Millisecond)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	time.Sleep(5 * time.Millisecond)

	if _, err := Parse(testSecret, token); err == nil {
		t.Fatal("Parse() expected error for expired token")
	}
}

func TestParseWrongSecret(t *testing.T) {
	token, err := Sign(testSecret, 1, "admin", time.Hour)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	other := "fedcba9876543210fedcba9876543210"
	if _, err := Parse(other, token); err == nil {
		t.Fatal("Parse() expected error for wrong secret")
	}
}

func TestParseEmptyToken(t *testing.T) {
	if _, err := Parse(testSecret, "   "); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Parse() error = %v, want %v", err, ErrInvalidToken)
	}
}

func TestMetadataRoundTrip(t *testing.T) {
	ctx := InjectUserID(context.Background(), 42)

	outgoing, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("outgoing metadata missing")
	}
	incoming := metadata.NewIncomingContext(context.Background(), outgoing)

	id, ok := UserIDFromIncoming(incoming)
	if !ok {
		t.Fatal("UserIDFromIncoming() did not find user id")
	}
	if id != 42 {
		t.Errorf("UserIDFromIncoming() = %d, want 42", id)
	}
}

func TestUserIDFromIncomingMissing(t *testing.T) {
	if _, ok := UserIDFromIncoming(context.Background()); ok {
		t.Fatal("UserIDFromIncoming() should fail without metadata")
	}
}
