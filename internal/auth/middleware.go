package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

// RequireAuth parses the Bearer JWT and injects the authenticated user's ID into
// the request context. Handlers must read the user ID from context, never from
// a client-supplied field, so a caller can never act on another user's behalf.
func RequireAuth(issuer *TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"missing bearer token"}`, http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(header, "Bearer ")

			userID, err := issuer.Parse(tokenString)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext reads the authenticated user's ID set by RequireAuth.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDContextKey).(int64)
	return id, ok
}
