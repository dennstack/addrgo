package main

import (
	"log"
	"net/http"

	"github.com/dennstack/addrgo/api"
	"github.com/dennstack/addrgo/middleware"
)

func StartApiServer() {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/search", api.SearcHandler)

	middlewareStack := middleware.CreateStack(
		middleware.LoggingMiddleware,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewareStack(router),
	}
	log.Println("Starting server on :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
