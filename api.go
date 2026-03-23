package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dennstack/addrgo/api"
	"github.com/dennstack/addrgo/middleware"
)

func StartApiServer() {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/search", api.SearchHandler)
	router.HandleFunc("POST /api/validate", api.ValidateHandler)

	middlewareStack := middleware.CreateStack(
		middleware.LoggingMiddleware,
	)

	server := http.Server{
		Addr:         ":8080",
		Handler:      middlewareStack(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped.")
}
