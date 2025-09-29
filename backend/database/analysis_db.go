package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"pikalink-backend/models"
	"pikalink-backend/utils"

	_ "modernc.org/sqlite"
)

func GetAnalysis(month, shortURL string) (*models.AnalysisResult, error) {
	// Use centralized path logic
	dbPath := utils.GetLogDbPath(month)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log database: %w", err)
	}
	defer db.Close()

	result := &models.AnalysisResult{
		StatusDistribution:  make(map[string]int64),
		UADistribution:      make(map[string]int64),
		RefererDistribution: make(map[string]int64),
	}

	// Get total clicks
	query := "SELECT COUNT(*) FROM access_logs"
	args := []interface{}{}
	if shortURL != "" {
		query += " WHERE short_url = ?"
		args = append(args, shortURL)
	}
	err = db.QueryRow(query, args...).Scan(&result.TotalClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get status distribution
	query = "SELECT http_status, COUNT(*) FROM access_logs"
	if shortURL != "" {
		query += " WHERE short_url = ?"
	}
	query += " GROUP BY http_status"
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get status distribution: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var status int
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			log.Printf("Error scanning status distribution: %v", err)
			continue
		}
		result.StatusDistribution[strconv.Itoa(status)] = count
	}

	// Get UA distribution
	query = "SELECT user_agent, COUNT(*) FROM access_logs"
	if shortURL != "" {
		query += " WHERE short_url = ?"
	}
	query += " GROUP BY user_agent ORDER BY COUNT(*) DESC"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get UA distribution: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var ua string
		var count int64
		if err := rows.Scan(&ua, &count); err != nil {
			log.Printf("Error scanning UA distribution: %v", err)
			continue
		}
		result.UADistribution[ua] = count
	}

	// Get referer distribution
	query = "SELECT referer, COUNT(*) FROM access_logs"
	if shortURL != "" {
		query += " WHERE short_url = ?"
	}
	query += " GROUP BY referer ORDER BY COUNT(*) DESC"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get referer distribution: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var referer string
		var count int64
		if err := rows.Scan(&referer, &count); err != nil {
			log.Printf("Error scanning referer distribution: %v", err)
			continue
		}
		result.RefererDistribution[referer] = count
	}

	// Get all short URLs
	rows, err = db.Query("SELECT DISTINCT short_url FROM access_logs")
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct short URLs: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			log.Printf("Error scanning short url: %v", err)
			continue
		}
		result.ShortURLs = append(result.ShortURLs, url)
	}

	return result, nil
}
