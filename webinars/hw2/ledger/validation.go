package ledger

import (
	"errors"
)

type Validatable interface {
	Validate() error
}

func (t Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("Размер транзакции должен быть положительным.")
	}
	if t.Category == "" {
		return errors.New("Не задана категория.")
	}
	if t.Date.IsZero() {
		return errors.New("Дата транзакции не задана.")
	}
	return nil
}

func (b Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("Бюджет должен быть положительным числом.")
	}
	if b.Category == "" {
		return errors.New("У бюджета должна быть категория.")
	}
	if b.Period == "" {
		return errors.New("У бюджета должен быть задан период.")
	}
	return nil
}

func CheckValid(v Validatable) error {
	return v.Validate()
}
