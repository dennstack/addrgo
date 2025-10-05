package main

import (
	"log"
	"net/http"

	"github.com/dennstack/addrgo/middleware"
)

func StartApiServer() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	middlewareStack := middleware.CreateStack(
		middleware.LoggingMiddleware,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewareStack(router),
	}
	log.Println("Starting server on :8080")
	server.ListenAndServe()
}
