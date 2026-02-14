package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aym-Aymen777/RSS-Aggregator/config"
	"github.com/Aym-Aymen777/RSS-Aggregator/handlers"
	"github.com/Aym-Aymen777/RSS-Aggregator/middleware"
	"github.com/Aym-Aymen777/RSS-Aggregator/services"
	"github.com/Aym-Aymen777/RSS-Aggregator/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	// Initialize MongoDB (single shared client)
	connectDB()

	jwtConfig := config.NewJWTConfig()
	tokenService := services.NewTokenService(jwtConfig)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1 := chi.NewRouter()
	//Health check endpoints
	v1.Get("/ready", handlerReadiness)
	v1.Get("/error", handlerErr)

	//CRUD operations endpoints for users
	v1.Post("/users/create", handlerCreateUser)
	v1.Post("/users/create-many", handlerCreateManyUsers)
	v1.Get("/users", handlerFindUserByEmail)
	v1.Put("/users/update", handlerUpdateUser)

	// Public routes (no authentication required)
	authCollection := MongoClient.Database("rssagg").Collection("auths")
	v1.Post("/auth/register", handlers.HandlerRagisterUser(authCollection))
	v1.Post("/auth/login", handlers.HandlerLoginUser(authCollection, tokenService))
	v1.Post("/auth/refresh", handlers.HandlerRefreshToken(authCollection, tokenService))
	// Protected routes (authentication required)
	postsCollection := MongoClient.Database("rssagg").Collection("posts")
	v1.Group(func(r chi.Router) {
		r.Use(middleware.AuthMidlleware(tokenService))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is a protected route"))
		})
		r.Get("/user/profile", func(w http.ResponseWriter, r *http.Request) {
			user, ok := middleware.GetUserFromContext(r.Context())
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
				return
			}
			utils.RespondWithJSON(w, http.StatusOK, user)
		})
		r.Post("/posts/create", handlers.HandlerCreatePost(postsCollection))
		r.Get("/posts", handlers.HandlerGetPosts(postsCollection))
		r.Get("/posts/{id}", handlers.HandlerGetPostByID(postsCollection))
		r.Put("/posts/{id}", handlers.HandlerUpdatePost(postsCollection))
		r.Delete("/posts/{id}", handlers.HandlerDeletePost(postsCollection))

	})
	router.Mount("/v1", v1)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	/* // Log registered routes
	log.Printf("🗺️  Registered Routes:")

	// First walk the main router
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*", "/", )
		log.Printf("[%s]: %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Printf("Main router logging err: %s\n", err.Error())
	}

	// Then walk the v1 sub-router with prefix
	log.Printf("\n📍 V1 Routes:")
	walkFuncV1 := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*", "/", )
		log.Printf("[%s]: /v1%s\n", method, route)
		return nil
	}
	if err := chi.Walk(v1, walkFuncV1); err != nil {
		log.Printf("V1 router logging err: %s\n", err.Error())
	} */
	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	if err := MongoClient.Disconnect(ctx); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	}

	log.Println("✅ Server stopped cleanly")
}
