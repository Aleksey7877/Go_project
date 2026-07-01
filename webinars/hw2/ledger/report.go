package ledger

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const reportSummaryCacheKey = "report:summary"

const reportSummaryCacheTTL = 5 * time.Minute

func buildReportSummaryFromDB() (ReportSummary, error) {
	database, err := requireDB()
	if err != nil {
		return ReportSummary{}, err
	}

	query := `
		select category, coalesce(sum(amount), 0) as total
		from expenses
		group by category
		order by category
	`

	rows, err := database.Query(query)
	if err != nil {
		return ReportSummary{}, err
	}
	defer rows.Close()

	summary := ReportSummary{
		Categories: make([]CategorySummary, 0),
	}

	for rows.Next() {
		var categorySummary CategorySummary

		err := rows.Scan(
			&categorySummary.Category,
			&categorySummary.Total,
		)
		if err != nil {
			return ReportSummary{}, err
		}

		summary.Total += categorySummary.Total
		summary.Categories = append(summary.Categories, categorySummary)
	}

	if err := rows.Err(); err != nil {
		return ReportSummary{}, err
	}

	return summary, nil
}

func GetReportSummary() (ReportSummary, error) {
	ctx := context.Background()

	cache, err := requireCache()
	if err == nil {
		cached, err := cache.Get(ctx, reportSummaryCacheKey).Result()
		if err == nil {
			var summary ReportSummary

			err = json.Unmarshal([]byte(cached), &summary)
			if err == nil {
				log.Println("report summary cache hit")
				return summary, nil
			}

			log.Printf("failed to unmarshal report summary cache: %v", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Printf("failed to get report summary cache: %v", err)
		}
	}

	summary, err := buildReportSummaryFromDB()
	if err != nil {
		return ReportSummary{}, err
	}

	cache, err = requireCache()
	if err == nil {
		data, err := json.Marshal(summary)
		if err == nil {
			err = cache.Set(ctx, reportSummaryCacheKey, data, reportSummaryCacheTTL).Err()
			if err != nil {
				log.Printf("failed to set report summary cache: %v", err)
			} else {
				log.Println("report summary cache miss")
			}
		}
	}

	return summary, nil
}
