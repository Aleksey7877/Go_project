package main

import (
	"fmt"
	"gateway/internal/api"
	"ledger"
	"log"
	"net/http"
)

func main() {
	if err := ledger.InitDB(); err != nil {
		log.Fatalf("failed to init ledger db: %v", err)
	}

	if err := ledger.InitCache(); err != nil {
		log.Fatalf("failed to init redis: %v", err)
	}
	defer ledger.CloseCache()
	
	defer ledger.CloseDB()
	mux := http.NewServeMux()

	api.RegisterRoutes(mux)

	handler := api.LoggingMiddleware(mux)

	fmt.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
