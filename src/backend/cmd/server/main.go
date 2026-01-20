// Package main provides the entry point for the backend server.
package main

import (
	"log"
	"net/http"
	"time"

	"backend/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	// Register handlers
	mux.Handle("/", handler.NewRootHandler())
	mux.Handle("/health", handler.NewHealthHandler())

	port := ":8080"
	log.Printf("Server starting on http://localhost%s\n", port)

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
