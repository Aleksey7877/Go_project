package ledger

import (
	"errors"
	"sync"
	"time"
)

var transactions []Transaction

var budgets map[string]map[string]Budget

var ErrBudgetExceeded = errors.New("budget exceeded")
var ErrBudgetNotCreated = errors.New("budget is not created")
var ErrBudgetPeriodWrong = errors.New("budget period must be a year")

var mu sync.RWMutex

func AddTransaction(tx Transaction) (Transaction, error) {

	err := CheckValid(tx)
	if err != nil {
		return Transaction{}, err
	}

	mu.Lock()
	defer mu.Unlock()

	sum := tx.Amount
	year := tx.Date.Format("2006")

	for _, existingTx := range transactions {
		if existingTx.Category == tx.Category &&
			existingTx.Date.Format("2006") == year {
			sum += existingTx.Amount
		}
	}
	yearBudget := tx.Date.Format("2006")

	budget, exists := budgets[yearBudget][tx.Category]

	if exists {
		if sum > budget.Limit {
			return Transaction{}, ErrBudgetExceeded
		}

		tx.ID = len(transactions) + 1
		transactions = append(transactions, tx)
		return tx, nil
	}
	return Transaction{}, ErrBudgetNotCreated
}

func ListTransactions() []Transaction {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]Transaction, len(transactions))
	copy(result, transactions)

	return result
}

func ListBudgets() []Budget {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]Budget, 0)

	for _, yearBudgets := range budgets {
		for _, budget := range yearBudgets {
			result = append(result, budget)
		}
	}

	return result
}

func SetBudget(budget Budget) (Budget, error) {
	err := CheckValid(budget)
	if err != nil {
		return Budget{}, err
	}

	mu.Lock()
	defer mu.Unlock()

	year := budget.Period
	_, err = time.Parse("2006", year)
	if err != nil {
		return Budget{}, ErrBudgetPeriodWrong
	}

	if budgets == nil {
		budgets = make(map[string]map[string]Budget)
		if budgets[year] == nil {
			budgets[year] = make(map[string]Budget)
		}
		budgets[year][budget.Category] = budget
		return budget, nil
	}
	if budgets[year] == nil {
		budgets[year] = make(map[string]Budget)
	}
	budgets[year][budget.Category] = budget
	return budget, nil
}
