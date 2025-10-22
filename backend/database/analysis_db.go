package database

import (
	"fmt"
	"log"
	"strconv"

	"pikalink-backend/models"
)

func GetAnalysis(month, shortURL, statusCode string) (*models.AnalysisResult, error) {
	// Use connection pool for better performance
	pool := GetAnalysisPool()
	db, err := pool.GetConnection(month)
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection for month %s: %w", month, err)
	}
	// Note: We don't defer db.Close() here because the pool manages the connection lifecycle

	result := &models.AnalysisResult{
		StatusDistribution:  make(map[string]int64),
		UADistribution:      make(map[string]int64),
		RefererDistribution: make(map[string]int64),
		BrowserDistribution: make(map[string]int64),
		OSDistribution:      make(map[string]int64),
		DeviceDistribution:  make(map[string]int64),
		BotDistribution:     make(map[string]int64),
	}

	// Build WHERE clause and args for filtering
	whereClause := ""
	args := []interface{}{}
	
	if shortURL != "" && statusCode != "" {
		whereClause = " WHERE short_url = ? AND http_status = ?"
		args = append(args, shortURL, statusCode)
	} else if shortURL != "" {
		whereClause = " WHERE short_url = ?"
		args = append(args, shortURL)
	} else if statusCode != "" {
		whereClause = " WHERE http_status = ?"
		args = append(args, statusCode)
	}

	// Get total clicks
	query := "SELECT COUNT(*) FROM access_logs" + whereClause
	err = db.QueryRow(query, args...).Scan(&result.TotalClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get status distribution
	query = "SELECT http_status, COUNT(*) FROM access_logs" + whereClause + " GROUP BY http_status ORDER BY COUNT(*) DESC LIMIT 101"
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get status distribution: %w", err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var status int
		var cnt int64
		if err := rows.Scan(&status, &cnt); err != nil {
			log.Printf("Error scanning status distribution: %v", err)
			continue
		}
		count++
		if count <= 100 {
			result.StatusDistribution[strconv.Itoa(status)] = cnt
		} else {
			// We found the 101st row, add "..." and break early
			result.StatusDistribution["... (only shows first 100 results)"] = 0
			break
		}
	}

	// Get UA distribution
	query = "SELECT user_agent, COUNT(*) FROM access_logs" + whereClause + " GROUP BY user_agent ORDER BY COUNT(*) DESC LIMIT 101"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get UA distribution: %w", err)
	}
	defer rows.Close()
	count = 0
	for rows.Next() {
		var ua string
		var cnt int64
		if err := rows.Scan(&ua, &cnt); err != nil {
			log.Printf("Error scanning UA distribution: %v", err)
			continue
		}
		count++
		if count <= 100 {
			result.UADistribution[ua] = cnt
		} else {
			// We found the 101st row, add "..." and break early
			result.UADistribution["... (only shows first 100 results)"] = 0
			break
		}
	}

	// Get referer distribution
	query = "SELECT referer, COUNT(*) FROM access_logs" + whereClause + " GROUP BY referer ORDER BY COUNT(*) DESC LIMIT 101"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get referer distribution: %w", err)
	}
	defer rows.Close()
	count = 0
	for rows.Next() {
		var referer string
		var cnt int64
		if err := rows.Scan(&referer, &cnt); err != nil {
			log.Printf("Error scanning referer distribution: %v", err)
			continue
		}
		count++
		if count <= 100 {
			result.RefererDistribution[referer] = cnt
		} else {
			// We found the 101st row, add "..." and break early
			result.RefererDistribution["... (only shows first 100 results)"] = 0
			break
		}
	}

	// Get browser distribution
	query = "SELECT COALESCE(browser, 'Unknown'), COUNT(*) FROM access_logs" + whereClause + " GROUP BY browser ORDER BY COUNT(*) DESC LIMIT 101"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser distribution: %w", err)
	}
	defer rows.Close()
	count = 0
	for rows.Next() {
		var browser string
		var cnt int64
		if err := rows.Scan(&browser, &cnt); err != nil {
			log.Printf("Error scanning browser distribution: %v", err)
			continue
		}
		count++
		if count <= 100 {
			result.BrowserDistribution[browser] = cnt
		} else {
			// We found the 101st row, add "..." and break early
			result.BrowserDistribution["... (only shows first 100 results)"] = 0
			break
		}
	}

	// Get OS distribution
	query = "SELECT COALESCE(os, 'Unknown'), COUNT(*) FROM access_logs" + whereClause + " GROUP BY os ORDER BY COUNT(*) DESC LIMIT 101"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get OS distribution: %w", err)
	}
	defer rows.Close()
	count = 0
	for rows.Next() {
		var os string
		var cnt int64
		if err := rows.Scan(&os, &cnt); err != nil {
			log.Printf("Error scanning OS distribution: %v", err)
			continue
		}
		count++
		if count <= 100 {
			result.OSDistribution[os] = cnt
		} else {
			// We found the 101st row, add "..." and break early
			result.OSDistribution["... (only shows first 100 results)"] = 0
			break
		}
	}

	// Get device type distribution
	query = "SELECT COALESCE(device_type, 'Unknown'), COUNT(*) FROM access_logs" + whereClause + " GROUP BY device_type ORDER BY COUNT(*) DESC LIMIT 101"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get device distribution: %w", err)
	}
	defer rows.Close()
	count = 0
	for rows.Next() {
		var device string
		var cnt int64
		if err := rows.Scan(&device, &cnt); err != nil {
			log.Printf("Error scanning device distribution: %v", err)
			continue
		}
		count++
		if count <= 100 {
			result.DeviceDistribution[device] = cnt
		} else {
			// We found the 101st row, add "..." and break early
			result.DeviceDistribution["... (only shows first 100 results)"] = 0
			break
		}
	}

	// Get bot distribution
	query = "SELECT CASE WHEN is_bot = 1 THEN 'Bot' ELSE 'Human' END, COUNT(*) FROM access_logs" + whereClause + " GROUP BY is_bot ORDER BY COUNT(*) DESC"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot distribution: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var botType string
		var cnt int64
		if err := rows.Scan(&botType, &cnt); err != nil {
			log.Printf("Error scanning bot distribution: %v", err)
			continue
		}
		result.BotDistribution[botType] = cnt
	}

	// Get all short URLs (only 200 and 302 status codes)
	rows, err = db.Query("SELECT DISTINCT short_url FROM access_logs WHERE http_status IN (200, 302)")
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

	// Get all distinct status codes
	rows, err = db.Query("SELECT DISTINCT http_status FROM access_logs ORDER BY http_status")
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct status codes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var status int
		if err := rows.Scan(&status); err != nil {
			log.Printf("Error scanning status code: %v", err)
			continue
		}
		result.StatusCodes = append(result.StatusCodes, strconv.Itoa(status))
	}

	return result, nil
}
