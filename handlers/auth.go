package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Aym-Aymen777/RSS-Aggregator/services"
	"github.com/Aym-Aymen777/RSS-Aggregator/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func HandlerRagisterUser(coll *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Parse the request body
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Here you would typically add logic to save the user to the database
		if req.Username == "" || req.Email == "" || req.Password == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}
		if len(req.Password) < 6 {
			http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
			return
		}
		err := services.RegisterUser(coll, req.Username, req.Email, req.Password)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		// Respond with a success message
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))

	}
}

func HandlerLoginUser(coll *mongo.Collection, tokenService *services.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Parse the request body
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}
		user, err := services.LoginUser(coll, req.Email, req.Password)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		userID := user.ID
		email := req.Email
		// Generate Token Pair
		tokens, err := tokenService.GenerateTokens(userID, email)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		// Set access token cookie (short-lived)
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokens.AccessToken,
			Path:     "/", 
			MaxAge:   15 * 60,                 // 15 minutes in seconds
			HttpOnly: true,                    // Prevents JavaScript access (XSS protection)
			Secure:   false,                   //TODO Only sent over HTTPS (set to false in development)
			SameSite: http.SameSiteStrictMode, // CSRF protection
		})

		// Set refresh token cookie (long-lived)
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Path:     "/", // Only sent to refresh endpoint
			MaxAge:   7 * 24 * 60 * 60, // 7 days in seconds
			HttpOnly: true,
			Secure:   false, //TODO Set to false in development
			SameSite: http.SameSiteStrictMode,
		})

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Login successful",
			"user":    user.Username,
			"tokens":  tokens,
		})
	}
}

func HandlerRefreshToken(coll *mongo.Collection, tokenService *services.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Get the refresh token from the cookie
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "Refresh token not provided", http.StatusUnauthorized)
			return
		}
		refreshToken := cookie.Value
		// Validate the refresh token and get the user ID
		claims, err := tokenService.ValidateRefreshToken(refreshToken)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}
		userID := claims.UserID
		email := claims.Email
		// Generate new token pair
		tokens, err := tokenService.GenerateTokens(userID, email)
		if err != nil {
			http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}
		// Set new access token cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "access_token",
			Value: tokens.AccessToken,

			Path:     "/",
			MaxAge:   15 * 60, // 15 minutes in seconds
			HttpOnly: true,
			Secure:   false, //TODO Only sent over HTTPS (set to false in development)
			SameSite: http.SameSiteStrictMode,
		})

		// Set new refresh token cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "refresh_token",
			Value: tokens.RefreshToken,

			Path:     "/v1/auth/refresh",    // Only sent to refresh endpoint
			MaxAge:   7 * 24 * 60 * 60, // 7 days in seconds
			HttpOnly: true,
			Secure:   false, //TODO Set to false in development
			SameSite: http.SameSiteStrictMode,
		})
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Token refreshed successfully",
			"tokens":  tokens,
		})
	}
}
