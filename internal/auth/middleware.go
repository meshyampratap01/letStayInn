package auth

import (
	"context"
	"net/http"
	"strings"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
)

// JWTAuthMiddleware protects endpoints with JWT authentication
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		// Store claims in context for downstream handlers
		ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, contextkeys.UserNameKey, claims.Username)
		ctx = context.WithValue(ctx, contextkeys.UserRoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
