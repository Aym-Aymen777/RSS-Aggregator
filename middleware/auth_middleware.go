package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/Aym-Aymen777/RSS-Aggregator/models"
	"github.com/Aym-Aymen777/RSS-Aggregator/services"
)

type contextKey string

const UserContextKey = contextKey("auth")

func AuthMidlleware(tokenService *services.TokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("🔒 AuthMiddleware: Checking authentication for %s", r.URL.Path)
			// Extract the token from the Authorization header
			var tokenString string
			cookie, err := r.Cookie("access_token")
			if err == nil {
				tokenString = cookie.Value
			} else {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
					tokenString = authHeader[7:]
				}
			}

			// Validate the token using the tokenService
			if tokenString == "" {
				log.Printf("🔒 AuthMiddleware: No token provided")
				http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
				return
			}
			claims, err := tokenService.ValidateAccessToken(tokenString)
			if err != nil {
				log.Printf("🔒 AuthMiddleware: Invalid token: %v", err)
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}
			// Store claims in request context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			r = r.WithContext(ctx)

			// If the token is valid, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user claims from request context
func GetUserFromContext(ctx context.Context) (*models.AccessTokenClaims, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.AccessTokenClaims)
	return user, ok
}
