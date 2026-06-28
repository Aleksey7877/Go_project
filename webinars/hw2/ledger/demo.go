package main

import "fmt"

func handleAdd(tx Transaction) {
	err := AddTransaction(tx)
	if err != nil {
		fmt.Println("Ошибка добавления транзакции:", err)
		return
	}
	fmt.Println("Успешно добавлена транзакция:", tx.Category, tx.Amount)
}

func handleBudget(b Budget) {
	err := SetBudget(b)
	if err != nil {
		fmt.Println("Ошибка создания бюджета:", err)
		return
	}
}
