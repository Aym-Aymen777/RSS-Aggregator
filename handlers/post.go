package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Aym-Aymen777/RSS-Aggregator/utils"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func HandlerCreatePost(coll *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// get the fields from the client
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Link        string `json:"link"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Title == "" || req.Description == "" || req.Link == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}
		coll.InsertOne(r.Context(), map[string]any{
			"title":       req.Title,
			"description": req.Description,
			"link":        req.Link,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		})
		// Respond with a success message
		utils.RespondWithJSON(w, http.StatusCreated, map[string]any{
			"message": "Post created successfully",
		})
	}
}

func HandlerGetPosts(coll *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//force get methode
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		//get all the posts from the database
		cursor, err := coll.Find(r.Context(), map[string]any{})
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		var posts []map[string]any
		if err = cursor.All(r.Context(), &posts); err != nil {
			http.Error(w, "Failed to decode posts", http.StatusInternalServerError)
			return
		}

		// Respond with the posts
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"posts": posts,
		})
	}
}

func HandlerGetPostByID(coll *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//force get methode
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		//get the id from the url
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Missing post ID", http.StatusBadRequest)
			return
		}
		//get the post from the database
		var post map[string]any
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		err = coll.FindOne(r.Context(), map[string]any{"_id": objectID}).Decode(&post)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Post not found", http.StatusNotFound)
			} else {
				log.Printf("Error fetching post: %v", err)
				http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
			}
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, post)
	}
}

func HandlerUpdatePost(coll *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//force put methode
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		//get the id from the url
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Missing post ID", http.StatusBadRequest)
			return
		}
	    ObjectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		//get the fields from the client
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Link        string `json:"link"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Title == "" && req.Description == "" && req.Link == "" {
			http.Error(w, "At least one field (title, description, or link) must be provided", http.StatusBadRequest)
			return
		}
		// Update the post in the database
		update := bson.M{}
		if req.Title != "" {
			update["title"] = req.Title
		}
		if req.Description != "" {
			update["description"] = req.Description
		}
		if req.Link != "" {
			update["link"] = req.Link
		}
		update["updated_at"] = time.Now()
		result, err := coll.UpdateOne(r.Context(), map[string]any{"_id": ObjectID}, bson.M{"$set": update})
		if err != nil {
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}
		if result.MatchedCount == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Post updated successfully",
		})
	}
}

func HandlerDeletePost(coll *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//force delete methode	
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		//get the id from the url
		id := chi.URLParam(r, "id")	
		if id == "" {
			http.Error(w, "Missing post ID", http.StatusBadRequest)
			return
		}
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		//delete the post from the database
		result, err := coll.DeleteOne(r.Context(), map[string]any{"_id": objectID})
		if err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}
		if result.DeletedCount == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Post deleted successfully",
		})
	}
}
