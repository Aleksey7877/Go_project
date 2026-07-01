package ledger

import "time"

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

type CategorySummary struct {
	Category string
	Total    float64
}

type ReportSummary struct {
	Total      float64
	Categories []CategorySummary
}