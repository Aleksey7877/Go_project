package main

import (
	"fmt"
	"gateway/internal/api"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	api.RegisterRoutes(mux)

	handler := api.LoggingMiddleware(mux)

	fmt.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
