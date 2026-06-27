package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

type Budget struct {
	Category string
	Limit    float64
	Period   string
}

var transactions []Transaction

var budgets map[string]map[string]Budget

func main() {
	handleBudget(Budget{Category: "Food", Limit: 5000, Period: "2026"})
	handleBudget(Budget{Category: "Transport", Limit: 2000, Period: "2026"})
	handleBudget(Budget{Category: "Entertainment", Limit: 3000, Period: "2026"})
	handleBudget(Budget{Category: "Entertainment", Limit: 4004, Period: "2026"})
	handleBudget(Budget{Category: "Food", Limit: 7000, Period: "2027"})

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
	handleAdd(Transaction{Amount: 50.0, Category: "Transport", Description: "Taxi", Date: time.Now()})
	handleAdd(Transaction{Amount: 0.0, Category: "Entertainment", Description: "Movie", Date: time.Now()})
	handleAdd(Transaction{Amount: 111.0, Category: "Entertainment", Description: "Example", Date: time.Now().AddDate(-1, 0, 0)})
	handleAdd(Transaction{Amount: 111.0, Category: "Entertainment", Description: "Example", Date: time.Now().AddDate(-2, 0, 0)})
	handleAdd(Transaction{Amount: 5000, Category: "Food", Description: "Extra meal", Date: time.Now()})
	handleAdd(Transaction{Amount: 500, Category: "Food", Description: "Extra meal", Date: time.Now()})

	fmt.Println(ListTransactions())
}

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
	fmt.Println("Бюджет успешно создан:", b.Category, b.Period)
}

func AddTransaction(tx Transaction) error {

	if tx.Amount == 0 {
		return fmt.Errorf("Значение транзакции не может быть равно нулю")
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
	if budget.Limit == 0 {
		return fmt.Errorf("Значение бюджета не может быть равно нулю")
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

func LoadBudgets(r io.Reader) error {
	var listBudgets []Budget

	err := json.NewDecoder(r).Decode(&listBudgets)
	if err != nil {
		return fmt.Errorf("Ошибка чтения бюджетов из JSON: %w", err)
	}

	for _, budget := range listBudgets {
		err := SetBudget(budget)
		if err != nil {
			return fmt.Errorf("ошибка установки бюджета %s %s: %w", budget.Category, budget.Period, err)
		}
	}
	return nil
}
