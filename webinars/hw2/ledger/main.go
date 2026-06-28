package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	handleBudget(Budget{Category: "Food", Limit: 5000, Period: "2026"})
	handleBudget(Budget{Category: "Transport", Limit: 2000, Period: "2026"})
	handleBudget(Budget{Category: "Entertainment", Limit: 3000, Period: "2026"})
	handleBudget(Budget{Category: "Entertainment", Limit: 4004, Period: "2026"})
	handleBudget(Budget{Category: "Food", Limit: 7000, Period: "2027"})
	handleBudget(Budget{Category: "", Limit: 7000, Period: "2027"})
	handleBudget(Budget{Category: "Food", Limit: -3, Period: "2027"})
	handleBudget(Budget{Category: "Food", Limit: 7000})
	handleBudget(Budget{Limit: 7000, Period: "2027"})

	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Println("Ошибка открытия файла бюджетов:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	err = LoadBudgets(reader)
	if err != nil {
		fmt.Println("Ошибка загрузки бюджетов:", err)
		return
	}

	for year, cat := range budgets {
		for _, budget := range cat {
			fmt.Printf("Year: %s, Category: %s, Limit: %.2f\n", year, budget.Category, budget.Limit)
		}
	}

	handleAdd(Transaction{Amount: 100.0, Category: "Food", Description: "Lunch", Date: time.Now()})
	handleAdd(Transaction{Amount: 250.0, Category: "Food", Description: "Breakfast on Gili", Date: time.Now()})
	handleAdd(Transaction{Amount: 250.0, Category: "Food", Description: "Breakfast on Gili 2"})
	handleAdd(Transaction{Amount: 250.0, Category: "", Description: "Breakfast on Gili 3", Date: time.Now()})
	handleAdd(Transaction{Amount: 50.0, Category: "Transport", Description: "Taxi", Date: time.Now()})
	handleAdd(Transaction{Amount: 0.0, Category: "Entertainment", Description: "Movie", Date: time.Now()})
	handleAdd(Transaction{Amount: 111.0, Category: "Entertainment", Description: "Example", Date: time.Now().AddDate(-1, 0, 0)})
	handleAdd(Transaction{Amount: 111.0, Category: "Entertainment", Description: "Example", Date: time.Now().AddDate(-2, 0, 0)})
	handleAdd(Transaction{Amount: 5000, Category: "Food", Description: "Extra meal", Date: time.Now()})
	handleAdd(Transaction{Amount: 500, Category: "Food", Description: "Extra meal", Date: time.Now()})

	fmt.Println(ListTransactions())
}