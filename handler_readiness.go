package main

import (
	"context"
	"net/http"
	"time"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := MongoClient.Ping(ctx, nil); err != nil {
		http.Error(w, "MongoDB not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}