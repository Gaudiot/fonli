package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"gaudiot.com/fonli/core/security/tokens"
)

type contextKey string

const UserIDKey contextKey = "userID"

type errorResponse struct {
	Error string `json:"error"`
}

func AuthMiddleware(ts tokens.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := ts.ParseAccessToken(tokenString)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired access token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}
