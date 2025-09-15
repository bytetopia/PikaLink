package logging

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// LogEntry represents a single log entry with all required information
type LogEntry struct {
	Date       string
	Time       string
	ShortURL   string
	TargetURL  string
	HTTPStatus int
	CallerIP   string
	UserAgent  string
	Referer    string
}

// EnsureLogsDirectory creates the logs directory if it doesn't exist
func EnsureLogsDirectory() error {
	logsDir := "./logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err := os.MkdirAll(logsDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create logs directory: %v", err)
		}
	}
	return nil
}

// GetLogFileName generates the log file name for the current month
func GetLogFileName() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d.log", now.Year(), now.Month())
}

// GetLogFilePath returns the full path to the current month's log file
func GetLogFilePath() string {
	return filepath.Join("./logs", GetLogFileName())
}

// FormatLogEntry formats a log entry into a readable string
func FormatLogEntry(entry LogEntry) string {
	return fmt.Sprintf("%s | %s | %s | %s | %d | %s | %s | %s\n",
		entry.Date,
		entry.Time,
		entry.ShortURL,
		entry.TargetURL,
		entry.HTTPStatus,
		entry.CallerIP,
		entry.UserAgent,
		entry.Referer,
	)
}

// WriteLogEntry writes a log entry to the current month's log file
func WriteLogEntry(entry LogEntry) error {
	// Ensure logs directory exists
	if err := EnsureLogsDirectory(); err != nil {
		return err
	}

	// Get log file path
	logFilePath := GetLogFilePath()

	// Open file in append mode, create if it doesn't exist
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	// Format and write the log entry
	logLine := FormatLogEntry(entry)
	_, err = file.WriteString(logLine)
	if err != nil {
		return fmt.Errorf("failed to write log entry: %v", err)
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

	// Extract Referer
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "-"
	}

	// Create log entry
	entry := LogEntry{
		Date:       now.Format("2006-01-02"),
		Time:       now.Format("15:04:05"),
		ShortURL:   shortURL,
		TargetURL:  targetURL,
		HTTPStatus: httpStatus,
		CallerIP:   clientIP,
		UserAgent:  userAgent,
		Referer:    referer,
	}

	return WriteLogEntry(entry)
}