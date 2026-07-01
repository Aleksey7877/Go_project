package api

import (
	"encoding/json"
	"ledger"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	err := ledger.InitDB()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	ledger.CloseDB()

	os.Exit(code)
}

func setupTestHandler() http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	return mux
}

func TestCreateBudgetHandler(t *testing.T) {
	ledger.Reset()
	t.Cleanup(ledger.Reset)

	handler := setupTestHandler()

	body := strings.NewReader(`{"category":"food","limit":1000,"period":"2026"}`)

	req := httptest.NewRequest(http.MethodPost, "/api/budgets", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type application/json; charset=utf-8, got %q", contentType)
	}

	var response BudgetResponse

	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Category != "food" {
		t.Errorf("expected category food, got %q", response.Category)
	}

	if response.Limit != 1000 {
		t.Errorf("expected limit 1000, got %v", response.Limit)
	}

	if response.Period != "2026" {
		t.Errorf("expected period 2026, got %q", response.Period)
	}
}

func TestListBudgetsHandler(t *testing.T) {
	ledger.Reset()
	t.Cleanup(ledger.Reset)

	handler := setupTestHandler()

	postBody := strings.NewReader(`{"category":"food","limit":1000,"period":"2026"}`)
	postReq := httptest.NewRequest(http.MethodPost, "/api/budgets", postBody)
	postReq.Header.Set("Content-Type", "application/json")

	postRR := httptest.NewRecorder()

	handler.ServeHTTP(postRR, postReq)

	if postRR.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, postRR.Code, postRR.Body.String())
	}

	invalidPostBody1 := strings.NewReader(`{"category":"food","limit":1000}`)
	invalidPostBody2 := strings.NewReader(`{"category":"food","limit":0,"period":"2026"}`)

	invalidPostReq1 := httptest.NewRequest(http.MethodPost, "/api/budgets", invalidPostBody1)
	invalidPostReq1.Header.Set("Content-Type", "application/json")

	invalidPostRR1 := httptest.NewRecorder()

	handler.ServeHTTP(invalidPostRR1, invalidPostReq1)

	if invalidPostRR1.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, invalidPostRR1.Code, invalidPostRR1.Body.String())
	}

	invalidPostReq2 := httptest.NewRequest(http.MethodPost, "/api/budgets", invalidPostBody2)
	invalidPostReq2.Header.Set("Content-Type", "application/json")

	invalidPostRR2 := httptest.NewRecorder()

	handler.ServeHTTP(invalidPostRR2, invalidPostReq2)

	if invalidPostRR2.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, invalidPostRR2.Code, invalidPostRR2.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/budgets", nil)
	getRR := httptest.NewRecorder()

	handler.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, getRR.Code, getRR.Body.String())
	}

	var responseGet []BudgetResponse

	err := json.NewDecoder(getRR.Body).Decode(&responseGet)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(responseGet) != 1 {
		t.Fatalf("expected len of body 1, got %d", len(responseGet))
	}

	budget := responseGet[0]

	if budget.Category != "food" {
		t.Errorf("expected category food, got %q", budget.Category)
	}

	if budget.Limit != 1000 {
		t.Errorf("expected limit 1000, got %v", budget.Limit)
	}

	if budget.Period != "2026" {
		t.Errorf("expected period 2026, got %q", budget.Period)
	}
}

func TestTransactionHandlers(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ledger.Reset()
		t.Cleanup(ledger.Reset)

		handler := setupTestHandler()

		body := strings.NewReader(`{"category":"food","limit":1000,"period":"2026"}`)

		req := httptest.NewRequest(http.MethodPost, "/api/budgets", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "application/json; charset=utf-8" {
			t.Errorf("expected Content-Type application/json; charset=utf-8, got %q", contentType)
		}

		var response BudgetResponse

		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Category != "food" {
			t.Errorf("expected category food, got %q", response.Category)
		}

		if response.Limit != 1000 {
			t.Errorf("expected limit 1000, got %v", response.Limit)
		}

		if response.Period != "2026" {
			t.Errorf("expected period 2026, got %q", response.Period)
		}

		bodyTransaction := strings.NewReader(`{"category":"food","amount":100,"date":"2026-09-10","description":"lunch"}`)

		reqTransaction := httptest.NewRequest(http.MethodPost, "/api/transactions", bodyTransaction)
		reqTransaction.Header.Set("Content-Type", "application/json")

		rrTransaction := httptest.NewRecorder()

		handler.ServeHTTP(rrTransaction, reqTransaction)

		if rrTransaction.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, rrTransaction.Code, rrTransaction.Body.String())
		}

		var createdTransaction TransactionResponse

		err = json.NewDecoder(rrTransaction.Body).Decode(&createdTransaction)
		if err != nil {
			t.Fatalf("failed to decode transaction response: %v", err)
		}

		if createdTransaction.Category != "food" {
			t.Errorf("expected category food, got %q", createdTransaction.Category)
		}

		if createdTransaction.Amount != 100 {
			t.Errorf("expected amount 100, got %v", createdTransaction.Amount)
		}

		if createdTransaction.Date != "2026-09-10" {
			t.Errorf("expected date 2026-09-10, got %q", createdTransaction.Date)
		}

		if createdTransaction.Description != "lunch" {
			t.Errorf("expected description lunch, got %q", createdTransaction.Description)
		}

		if createdTransaction.ID != 1 {
			t.Errorf("expected ID 1, got %d", createdTransaction.ID)
		}

		getReq := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
		getRR := httptest.NewRecorder()

		handler.ServeHTTP(getRR, getReq)

		if getRR.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, getRR.Code, getRR.Body.String())
		}

		var responseGet []TransactionResponse

		errTr := json.NewDecoder(getRR.Body).Decode(&responseGet)
		if errTr != nil {
			t.Fatalf("failed to decode response: %v", errTr)
		}

		if len(responseGet) != 1 {
			t.Fatalf("expected len of body 1, got %d", len(responseGet))
		}

		transaction := responseGet[0]

		if transaction.Category != "food" {
			t.Errorf("expected category food, got %q", transaction.Category)
		}

		if transaction.Amount != 100 {
			t.Errorf("expected amount 100, got %v", transaction.Amount)
		}

		if transaction.Date != "2026-09-10" {
			t.Errorf("expected date 2026-09-10, got %q", transaction.Date)
		}

		if transaction.Description != "lunch" {
			t.Errorf("expected description 2026, got %q", transaction.Description)
		}

		if transaction.ID != 1 {
			t.Errorf("expected ID 1, got %d", transaction.ID)
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		ledger.Reset()
		t.Cleanup(ledger.Reset)

		handler := setupTestHandler()

		InvalidBodyTransaction := strings.NewReader(`{"category":"food","amount":100,"description":"lunch"`)

		InvalidReqTransaction := httptest.NewRequest(http.MethodPost, "/api/transactions", InvalidBodyTransaction)
		InvalidReqTransaction.Header.Set("Content-Type", "application/json")

		InvalidRrTransaction := httptest.NewRecorder()

		handler.ServeHTTP(InvalidRrTransaction, InvalidReqTransaction)

		if InvalidRrTransaction.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, InvalidRrTransaction.Code, InvalidRrTransaction.Body.String())
		}
	})

	t.Run("exceeded", func(t *testing.T) {
		ledger.Reset()
		t.Cleanup(ledger.Reset)

		handler := setupTestHandler()

		// 1. Создать бюджет
		budgetBody := strings.NewReader(`{"category":"food","limit":1000,"period":"2026"}`)
		budgetReq := httptest.NewRequest(http.MethodPost, "/api/budgets", budgetBody)
		budgetReq.Header.Set("Content-Type", "application/json")

		budgetRR := httptest.NewRecorder()
		handler.ServeHTTP(budgetRR, budgetReq)

		if budgetRR.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, budgetRR.Code, budgetRR.Body.String())
		}

		// 2. Отправить превышающую транзакцию
		exceedBody := strings.NewReader(`{"category":"food","amount":1000.1,"date":"2026-09-10","description":"lunch"}`)
		exceedReq := httptest.NewRequest(http.MethodPost, "/api/transactions", exceedBody)
		exceedReq.Header.Set("Content-Type", "application/json")

		exceedRR := httptest.NewRecorder()
		handler.ServeHTTP(exceedRR, exceedReq)

		if exceedRR.Code != http.StatusConflict {
			t.Fatalf("expected status %d, got %d, body: %s", http.StatusConflict, exceedRR.Code, exceedRR.Body.String())
		}

		var errorResponse ErrorResponse

		err := json.NewDecoder(exceedRR.Body).Decode(&errorResponse)
		if err != nil {
			t.Fatalf("failed to decode error response: %v", err)
		}

		if errorResponse.Error != "budget exceeded" {
			t.Errorf("expected error budget exceeded, got %q", errorResponse.Error)
		}
	})
}
