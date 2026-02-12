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

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Login successful",
			"user":    user.Username,
			"tokens":  tokens,
		})
	}
}
