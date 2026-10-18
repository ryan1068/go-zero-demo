// Package auth provides JWT signing, parsing, HTTP authentication middleware,
// and helpers for propagating the authenticated user through contexts.
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/metadata"
)

const (
	// MetadataUserIDKey is the gRPC metadata key carrying the authenticated user ID.
	MetadataUserIDKey = "gateway-user-id"
	// MetadataUserRoleKey is the gRPC metadata key carrying the authenticated user role.
	MetadataUserRoleKey = "gateway-user-role"

	metadataHeaderUserID   = "Grpc-Metadata-User-Id"
	metadataHeaderUsername = "Grpc-Metadata-Username"
	metadataHeaderUserRole = "Grpc-Metadata-User-Role"
	minSecretLength        = 32
	tokenIssuer            = "go-zero-demo"
)

// RoleAdmin is the role value reserved for administrators.
const RoleAdmin int32 = 1

// RoleOperator is the role value reserved for day-to-day operators.
const RoleOperator int32 = 2

// IsAdmin reports whether the role can access the management backend.
func IsAdmin(role int32) bool {
	return role == RoleAdmin || role == RoleOperator
}

// IsSuperAdmin reports whether the role can perform destructive operations.
func IsSuperAdmin(role int32) bool {
	return role == RoleAdmin
}

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
	roleKey
)

// Claims is the JWT payload issued by the user service.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     int32  `json:"role"`
	jwt.RegisteredClaims
}

// ValidateSecret checks that the JWT signing secret is usable.
func ValidateSecret(secret string) error {
	if len(strings.TrimSpace(secret)) < minSecretLength {
		return ErrInvalidSecret
	}
	return nil
}

// Sign issues an HS256 JWT for the given user with a normal-user role.
func Sign(secret string, userID int64, username string, ttl time.Duration) (string, error) {
	return SignWithRole(secret, userID, username, 0, ttl)
}

// SignWithRole issues an HS256 JWT carrying the given role.
func SignWithRole(secret string, userID int64, username string, role int32, ttl time.Duration) (string, error) {
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
		Role:     role,
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

// WithUser returns a context carrying the authenticated identity.
func WithUser(ctx context.Context, userID int64, username string) context.Context {
	return WithUserRole(ctx, userID, username, 0)
}

// WithUserRole returns a context carrying the authenticated identity and role.
func WithUserRole(ctx context.Context, userID int64, username string, role int32) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, usernameKey, username)
	return context.WithValue(ctx, roleKey, role)
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

// Role returns the authenticated user role stored in ctx.
func Role(ctx context.Context) (int32, bool) {
	value := ctx.Value(roleKey)
	if value == nil {
		return 0, false
	}
	role, ok := value.(int32)
	return role, ok
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

// BearerToken returns the token carried by the Authorization header, if any.
func BearerToken(r *http.Request) string {
	token, _ := tokenFromRequest(r)
	return token
}

// TokenHash returns a stable SHA-256 fingerprint used for token blacklists.
func TokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// BlacklistKey returns the Redis key for a revoked token.
func BlacklistKey(token string) string {
	return "jwt:blacklist:" + TokenHash(token)
}
