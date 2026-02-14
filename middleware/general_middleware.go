package middleware

import (
	"log"
	"net/http"
)

func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Code BEFORE the handler runs
		log.Println("Before handler ⏳")

		// Call the next handler
		next.ServeHTTP(w, r)

		// Code AFTER the handler runs
		log.Println("After handler ⏳")
	})
}
