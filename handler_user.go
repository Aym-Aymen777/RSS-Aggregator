package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Aym-Aymen777/RSS-Aggregator/models"
	 "github.com/Aym-Aymen777/RSS-Aggregator/utils"
)

func handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	//force method post
	if r.Method != "POST" {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if user.Name == "" || user.Email == "" || user.Age <= 0 {
		respondError(w, http.StatusBadRequest, "Missing or invalid user fields")
		return
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err = utils.InsertUser(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func handlerCreateManyUsers(w http.ResponseWriter, r *http.Request) {
	//force method post
	if r.Method != "POST" {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var users []models.User
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = utils.InsertMany(users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create users")
		return
	}
	respondWithJSON(w, http.StatusCreated, users)
}

func handlerFindUserByEmail(w http.ResponseWriter, r *http.Request) {
	//force method get
	if r.Method != "GET" {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	email := r.URL.Query().Get("email")
	if email == "" {
		respondError(w, http.StatusBadRequest, "Email query parameter is required")
		return
	}
	results := utils.FindByQuery("email", email)
	if len(results) == 0 {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	respondWithJSON(w, http.StatusOK, results)
}

func handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	//force method put
	if r.Method != "PUT" {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "ID query parameter is required")
		return
	}

	utils.UpdateUser(id)
	respondWithJSON(w, http.StatusOK, "User updated successfully ✅")
}
