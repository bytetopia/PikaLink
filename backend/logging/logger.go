package logging

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"pikalink-backend/utils"
	_ "modernc.org/sqlite"
)

// LogEntry represents a single log entry with all required information
type LogEntry struct {
	Date           string
	Time           string
	ShortURL       string
	TargetURL      string
	HTTPStatus     int
	CallerIP       string
	UserAgent      string
	Referer        string
	Browser        string
	BrowserVersion string
	OS             string
	DeviceType     string
	IsBot          bool
}

// EnsureLogsDirectory creates the logs directory if it doesn't exist
func EnsureLogsDirectory() error {
	logsDir := utils.GetLogsDir()
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err := os.MkdirAll(logsDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create logs directory: %v", err)
		}
	}
	return nil
}

// GetLogDbFileName generates the log database file name for the current month
func GetLogDbFileName() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d.db", now.Year(), now.Month())
}

// GetLogDbFilePath returns the full path to the current month's log database file
func GetLogDbFilePath() string {
	now := time.Now()
	month := fmt.Sprintf("%d-%02d", now.Year(), now.Month())
	return utils.GetLogDbPath(month)
}

// initLogDB initializes the SQLite database for logging
func initLogDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to log database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping log database: %v", err)
	}

	// Optimize log database for write-heavy workload
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		// WAL mode improves concurrent write performance
	}
	if _, err := db.Exec("PRAGMA synchronous=NORMAL"); err != nil {
		// Reduce fsync calls for better performance
	}
	if _, err := db.Exec("PRAGMA cache_size=2000"); err != nil {
		// Larger cache for better performance
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS access_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT,
        time TEXT,
        short_url TEXT,
        target_url TEXT,
        http_status INTEGER,
        caller_ip TEXT,
        user_agent TEXT,
        referer TEXT,
        browser TEXT,
        browser_version TEXT,
        os TEXT,
        device_type TEXT,
        is_bot BOOLEAN
    );`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create access_logs table: %v", err)
	}

	return db, nil
}

// WriteLogEntry writes a log entry to the current month's log database
func WriteLogEntry(entry LogEntry) error {
	// Ensure logs directory exists
	if err := EnsureLogsDirectory(); err != nil {
		return err
	}

	// Get log db file path using centralized utils
	logDbPath := GetLogDbFilePath()

	// Initialize database and create table if not exists
	db, err := initLogDB(logDbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert log entry
	insertSQL := `INSERT INTO access_logs (date, time, short_url, target_url, http_status, caller_ip, user_agent, referer, browser, browser_version, os, device_type, is_bot)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(insertSQL, entry.Date, entry.Time, entry.ShortURL, entry.TargetURL, entry.HTTPStatus, entry.CallerIP, entry.UserAgent, entry.Referer, entry.Browser, entry.BrowserVersion, entry.OS, entry.DeviceType, entry.IsBot)
	if err != nil {
		return fmt.Errorf("failed to insert log entry into database: %v", err)
	}

	return nil
}

// LogLinkAccess logs a link access with information extracted from HTTP request
func LogLinkAccess(r *http.Request, shortURL, targetURL string, httpStatus int) error {
	now := time.Now()

	// Extract client IP
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Real-IP")
	}
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	// Extract User-Agent
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "-"
	}

	// Parse User-Agent for detailed information
	uaInfo := ParseUserAgent(userAgent)

	// Extract Referer
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "-"
	}

	// Create log entry
	entry := LogEntry{
		Date:           now.Format("2006-01-02"),
		Time:           now.Format("15:04:05"),
		ShortURL:       shortURL,
		TargetURL:      targetURL,
		HTTPStatus:     httpStatus,
		CallerIP:       clientIP,
		UserAgent:      userAgent,
		Referer:        referer,
		Browser:        uaInfo.Browser,
		BrowserVersion: uaInfo.BrowserVersion,
		OS:             uaInfo.OS,
		DeviceType:     uaInfo.DeviceType,
		IsBot:          uaInfo.IsBot,
	}

	return WriteLogEntry(entry)
}