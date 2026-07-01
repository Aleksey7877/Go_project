package api

import (
	"encoding/json"
	"errors"
	"ledger"
	"net/http"
)

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	tx, err := transactionFromRequest(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format")
		return
	}

	err = tx.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdTx, err := ledger.AddTransaction(tx)
	if err != nil {
		if errors.Is(err, ledger.ErrBudgetExceeded) {
			writeError(w, http.StatusConflict, "budget exceeded")
			return
		}

		if errors.Is(err, ledger.ErrBudgetNotCreated) {
			writeError(w, http.StatusConflict, "budget is not created")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := transactionToResponse(createdTx)
	writeJSON(w, http.StatusCreated, response)
}

func createBudgetHandler(w http.ResponseWriter, r *http.Request) {
	var budget CreateBudgetRequest

	err := json.NewDecoder(r.Body).Decode(&budget)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	b := budgetFromRequest(budget)

	err = b.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdBudget, err := ledger.SetBudget(b)
	if err != nil {
		if errors.Is(err, ledger.ErrBudgetPeriodWrong) {
			writeError(w, http.StatusBadRequest, "budget period must be a year")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := budgetToResponse(createdBudget)
	writeJSON(w, http.StatusCreated, response)
}

func listTransactionsHandler(w http.ResponseWriter, _ *http.Request) {
	transactions, err := ledger.ListTransactions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]TransactionResponse, 0, len(transactions))

	for _, tx := range transactions {
		response = append(response, transactionToResponse(tx))
	}

	writeJSON(w, http.StatusOK, response)
}

func listBudgetsHandler(w http.ResponseWriter, _ *http.Request) {
	budgets, err := ledger.ListBudgets()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]BudgetResponse, 0, len(budgets))

	for _, b := range budgets {
		response = append(response, budgetToResponse(b))
	}

	writeJSON(w, http.StatusOK, response)
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTransactionHandler(w, r)
	case http.MethodGet:
		listTransactionsHandler(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func budgetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createBudgetHandler(w, r)
	case http.MethodGet:
		listBudgetsHandler(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func reportSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	summary, err := ledger.GetReportSummary()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := reportSummaryToResponse(summary)
	writeJSON(w, http.StatusOK, response)
}