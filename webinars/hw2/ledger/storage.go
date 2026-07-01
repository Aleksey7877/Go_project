package ledger

import (
	"database/sql"
	"errors"
	"time"
)

var ErrBudgetExceeded = errors.New("budget exceeded")
var ErrBudgetNotCreated = errors.New("budget is not created")
var ErrBudgetPeriodWrong = errors.New("budget period must be a year")

func AddTransaction(tx Transaction) (Transaction, error) {
	err := CheckValid(tx)
	if err != nil {
		return Transaction{}, err
	}

	database, err := requireDB()
	if err != nil {
		return Transaction{}, err
	}

	year := tx.Date.Format("2006")

	var budgetLimit float64

	err = database.QueryRow(`
		select limit_amount
		from budgets
		where category = $1 AND period = $2
	`, tx.Category, year).Scan(&budgetLimit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Transaction{}, ErrBudgetNotCreated
		}

		return Transaction{}, err
	}

	start := time.Date(tx.Date.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	var currentSum float64

	err = database.QueryRow(`
		select coalesce(sum(amount), 0)
		from expenses
		where category = $1
		  and date >= $2
		  and date < $3
	`, tx.Category, start, end).Scan(&currentSum)
	if err != nil {
		return Transaction{}, err
	}

	if currentSum+tx.Amount > budgetLimit {
		return Transaction{}, ErrBudgetExceeded
	}

	err = database.QueryRow(`
		insert into expenses(amount, category, description, date)
		values($1, $2, $3, $4)
		returning id
	`, tx.Amount, tx.Category, tx.Description, tx.Date).Scan(&tx.ID)
	if err != nil {
		return Transaction{}, err
	}

	invalidateReportSummaryCache()

	return tx, nil
}

func ListTransactions() ([]Transaction, error) {

	database, err := requireDB()
	if err != nil {
		return nil, err
	}

	query := `
		select id, amount, category, coalesce(description, ''), date from expenses order by date desc, id desc;
	`

	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(&transaction.ID, &transaction.Amount, &transaction.Category, &transaction.Description, &transaction.Date)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func ListBudgets() ([]Budget, error) {
	database, err := requireDB()
	if err != nil {
		return nil, err
	}

	query := `
		select category, limit_amount, period from budgets order by period desc, category;
	`

	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgets := make([]Budget, 0)
	for rows.Next() {
		var budget Budget

		err := rows.Scan(&budget.Category, &budget.Limit, &budget.Period)
		if err != nil {
			return nil, err
		}

		budgets = append(budgets, budget)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return budgets, nil
}

func SetBudget(budget Budget) (Budget, error) {
	err := CheckValid(budget)
	if err != nil {
		return Budget{}, err
	}

	_, err = time.Parse("2006", budget.Period)
	if err != nil {
		return Budget{}, ErrBudgetPeriodWrong
	}

	database, err := requireDB()
	if err != nil {
		return Budget{}, err
	}

	query := `
		insert into budgets (category, limit_amount, period) 
		values ($1, $2, $3) 
		on conflict (category, period) 
		do update set limit_amount = excluded.limit_amount 
		returning category, limit_amount, period
		`

	var created Budget

	err = database.QueryRow(
		query,
		budget.Category,
		budget.Limit,
		budget.Period,
	).Scan(
		&created.Category,
		&created.Limit,
		&created.Period,
	)

	if err != nil {
		return Budget{}, err
	}

	return created, nil
}

func requireDB() (*sql.DB, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	return db, nil
}

func Reset() {
	if db == nil {
		return
	}

	_, _ = db.Exec(`
		TRUNCATE TABLE expenses, budgets RESTART IDENTITY
	`)
}