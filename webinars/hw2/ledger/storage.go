package main

import "fmt"

var transactions []Transaction

var budgets map[string]map[string]Budget

func AddTransaction(tx Transaction) error {

	err := CheckValid(tx)
	if err != nil {
		return err
	}
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
			return fmt.Errorf("Превышен лимит бюджета для категории: %s", tx.Category)
		}

		tx.ID = len(transactions) + 1
		transactions = append(transactions, tx)
		return nil
	}
	return fmt.Errorf("бюджет для года %s в категории %s не предусмотрен", year, tx.Category)
}

func ListTransactions() []Transaction {
	return transactions
}

func SetBudget(budget Budget) error {
	err := CheckValid(budget)
	if err != nil {
		return err
	}
	year := budget.Period

	if budgets == nil {
		budgets = make(map[string]map[string]Budget)
		if budgets[year] == nil {
			budgets[year] = make(map[string]Budget)
		}
		budgets[year][budget.Category] = budget
		fmt.Println("Бюджет успешно обновлен: ", budget.Category, budget.Period)
		return nil
	}
	if budgets[year] == nil {
		budgets[year] = make(map[string]Budget)
	}
	budgets[year][budget.Category] = budget
	fmt.Println("Бюджет успешно обновлен: ", budget.Category, budget.Period)
	return nil
}
