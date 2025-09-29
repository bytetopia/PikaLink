package models

type AnalysisResult struct {
	TotalClicks         int64             `json:"total_clicks"`
	StatusDistribution  map[string]int64  `json:"status_distribution"`
	UADistribution      map[string]int64  `json:"ua_distribution"`
	RefererDistribution map[string]int64  `json:"referer_distribution"`
	ShortURLs           []string          `json:"short_urls"`
}
