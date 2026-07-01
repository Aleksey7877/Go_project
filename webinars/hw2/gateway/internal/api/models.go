package api

type CategorySummaryResponse struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type ReportSummaryResponse struct {
	Total      float64                   `json:"total"`
	ByCategory []CategorySummaryResponse `json:"by_category"`
}