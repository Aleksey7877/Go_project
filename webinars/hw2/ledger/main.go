package main

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID     int
	Amount float64
	Category string
	Description string
	Date  time.Time
}

var transactions []Transaction

func main() {
	AddTransaction(Transaction{Amount: 100.0, Category: "Food", Description: "Lunch", Date: time.Now()})
	AddTransaction(Transaction{Amount: -250.0, Category: "Food", Description: "Breakfast on Gili", Date: time.Now()})
	AddTransaction(Transaction{Amount: 50.0, Category: "Transport", Description: "Taxi", Date: time.Now()})
	AddTransaction(Transaction{Amount: 0.0, Category: "Entertainment", Description: "Movie", Date: time.Now()})

	fmt.Println(ListTransactions())
}

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		fmt.Println("Значение транзакции не может быть равно нулю, номер операции:", len(transactions)+1)
		return fmt.Errorf("Значение транзакции не может быть равно нулю")
	}
	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)
	return nil
}

func ListTransactions() []Transaction {
	return transactions
}