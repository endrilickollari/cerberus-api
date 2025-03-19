package handlers

import (
	"context"
	"net/http"
	"strings"

	"remote-server-api/internal/api/response"
	"remote-server-api/internal/domain/auth"
)

// ContextKey is a type for context keys
type ContextKey string

// Context keys
const (
	UserIDKey    ContextKey = "userID"
	SessionIDKey ContextKey = "sessionID"
)

// AuthMiddleware handles authentication for protected routes
type AuthMiddleware struct {
	authService auth.Service
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate validates the JWT token and extracts claims
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token and extract claims
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			response.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Username)
		ctx = context.WithValue(ctx, SessionIDKey, claims.SessionID)

		// Call the next handler with the enhanced context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
