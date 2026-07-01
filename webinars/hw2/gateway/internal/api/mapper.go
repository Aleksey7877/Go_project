package api

import (
	"ledger"
	"time"
)

func transactionFromRequest(req CreateTransactionRequest) (ledger.Transaction, error) {
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return ledger.Transaction{}, err
	}
	tx := ledger.Transaction{
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Date:        parsedDate,
	}
	return tx, nil
}

func budgetFromRequest(req CreateBudgetRequest) ledger.Budget {
	b := ledger.Budget{
		Category: req.Category,
		Limit:    req.Limit,
		Period:   req.Period,
	}
	return b
}

func transactionToResponse(tx ledger.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date.Format("2006-01-02"),
	}
}

func budgetToResponse(b ledger.Budget) BudgetResponse {
	return BudgetResponse{
		Category: b.Category,
		Limit:    b.Limit,
		Period:   b.Period,
	}
}

func reportSummaryToResponse(summary ledger.ReportSummary) ReportSummaryResponse {
	response := ReportSummaryResponse{
		Total:      summary.Total,
		ByCategory: make([]CategorySummaryResponse, 0, len(summary.Categories)),
	}

	for _, category := range summary.Categories {
		response.ByCategory = append(response.ByCategory, CategorySummaryResponse{
			Category: category.Category,
			Total:    category.Total,
		})
	}

	return response
}
