package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/transactions", transactionHandler)
	mux.HandleFunc("/api/budgets", budgetHandler)
}