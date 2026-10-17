// Package auth provides JWT signing, parsing, HTTP authentication middleware,
// and helpers for propagating the authenticated user through contexts.
package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/metadata"
)

const (
	// MetadataUserIDKey is the gRPC metadata key carrying the authenticated user ID.
	MetadataUserIDKey = "gateway-user-id"

	metadataHeaderUserID   = "Grpc-Metadata-User-Id"
	metadataHeaderUsername = "Grpc-Metadata-Username"
	minSecretLength        = 32
	tokenIssuer            = "go-zero-demo"
)

var (
	// ErrInvalidSecret indicates the JWT secret is missing or too short.
	ErrInvalidSecret = errors.New("jwt secret must be at least 32 characters")
	// ErrInvalidToken indicates the token is missing, malformed, or expired.
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidUserID indicates the token does not carry a valid user ID.
	ErrInvalidUserID = errors.New("invalid user id")
	// ErrInvalidTTL indicates the requested token lifetime is not positive.
	ErrInvalidTTL = errors.New("invalid token ttl")
)

type contextKey int

const (
	userIDKey contextKey = iota
	usernameKey
)

// Claims is the JWT payload issued by the user service.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ValidateSecret checks that the JWT signing secret is usable.
func ValidateSecret(secret string) error {
	if len(strings.TrimSpace(secret)) < minSecretLength {
		return ErrInvalidSecret
	}
	return nil
}

// Sign issues an HS256 JWT for the given user.
func Sign(secret string, userID int64, username string, ttl time.Duration) (string, error) {
	if err := ValidateSecret(secret); err != nil {
		return "", err
	}
	if userID <= 0 {
		return "", ErrInvalidUserID
	}
	if ttl <= 0 {
		return "", ErrInvalidTTL
	}

	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Parse validates a JWT and returns its claims.
func Parse(secret, tokenString string) (*Claims, error) {
	if err := ValidateSecret(secret); err != nil {
		return nil, err
	}
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.UserID <= 0 {
		return nil, ErrInvalidUserID
	}
	return claims, nil
}

// UserID returns the authenticated user ID stored in ctx.
func UserID(ctx context.Context) (int64, bool) {
	value := ctx.Value(userIDKey)
	if value == nil {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok && id > 0
}

// Username returns the authenticated username stored in ctx.
func Username(ctx context.Context) (string, bool) {
	value := ctx.Value(usernameKey)
	if value == nil {
		return "", false
	}
	name, ok := value.(string)
	return name, ok && name != ""
}

// InjectUserID returns a context carrying userID as outgoing gRPC metadata.
func InjectUserID(ctx context.Context, userID int64) context.Context {
	return metadata.AppendToOutgoingContext(ctx, MetadataUserIDKey, strconv.FormatInt(userID, 10))
}

// UserIDFromIncoming returns the user ID carried by incoming gRPC metadata.
func UserIDFromIncoming(ctx context.Context) (int64, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, false
	}
	values := md.Get(MetadataUserIDKey)
	if len(values) == 0 {
		return 0, false
	}
	id, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
