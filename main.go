package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aym-Aymen777/RSS-Aggregator/handlers"
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
	v1.Put("/users/update",handlerUpdateUser)

	//Auth endpoints
	authCollection := MongoClient.Database("rssagg").Collection("auths")
	v1.Post("/auth/register", handlers.HandlerRagisterUser(authCollection))


	router.Mount("/v1", v1)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server
	go func() {
		log.Printf("🚀 Server running on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	MongoClient.Disconnect(ctx)

	log.Println("✅ Server stopped cleanly")
}
