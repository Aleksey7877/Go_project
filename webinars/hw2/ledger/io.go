package ledger

import (
	"encoding/json"
	"fmt"
	"io"
)

func LoadBudgets(r io.Reader) error {
	var listBudgets []Budget

	err := json.NewDecoder(r).Decode(&listBudgets)
	if err != nil {
		return fmt.Errorf("Ошибка чтения бюджетов из JSON: %w", err)
	}

	for _, budget := range listBudgets {
		_, err := SetBudget(budget)
		if err != nil {
			return fmt.Errorf("ошибка установки бюджета %s %s: %w", budget.Category, budget.Period, err)
		}
	}
	return nil
}
