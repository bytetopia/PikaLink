package models

type AnalysisResult struct {
	TotalClicks          int64             `json:"total_clicks"`
	StatusDistribution   map[string]int64  `json:"status_distribution"`
	URLDistribution      map[string]int64  `json:"url_distribution"`
	UADistribution       map[string]int64  `json:"ua_distribution"`
	RefererDistribution  map[string]int64  `json:"referer_distribution"`
	BrowserDistribution  map[string]int64  `json:"browser_distribution"`
	OSDistribution       map[string]int64  `json:"os_distribution"`
	DeviceDistribution   map[string]int64  `json:"device_distribution"`
	BotDistribution      map[string]int64  `json:"bot_distribution"`
	CountryDistribution  map[string]int64  `json:"country_distribution"`
	CityDistribution     map[string]int64  `json:"city_distribution"`
	ASNDistribution      map[string]int64  `json:"asn_distribution"`
	ShortURLs            []string          `json:"short_urls"`
	StatusCodes          []string          `json:"status_codes"`
}
