package logging

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"pikalink-backend/utils"
	_ "modernc.org/sqlite"
)

// AsyncLogger handles asynchronous logging operations
type AsyncLogger struct {
	logChannel   chan LogEntry
	db           *sql.DB
	dbMutex      sync.RWMutex
	currentMonth string
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	batchSize    int
	flushInterval time.Duration
}

var (
	asyncLogger *AsyncLogger
	loggerOnce  sync.Once
)

// GetAsyncLogger returns the singleton async logger instance
func GetAsyncLogger() *AsyncLogger {
	loggerOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		asyncLogger = &AsyncLogger{
			logChannel:    make(chan LogEntry, 1000), // Buffer for 1000 log entries
			ctx:           ctx,
			cancel:        cancel,
			batchSize:     50,                      // Process logs in batches of 50
			flushInterval: 5 * time.Second,         // Flush every 5 seconds
		}
		asyncLogger.start()
	})
	return asyncLogger
}

// start initializes the async logger and starts the worker goroutine
func (al *AsyncLogger) start() {
	// Initialize database connection
	if err := al.initializeDB(); err != nil {
		log.Printf("Failed to initialize async logger database: %v", err)
		return
	}

	al.wg.Add(1)
	go al.worker()

	// Start a goroutine to handle database rotation
	al.wg.Add(1)
	go al.dbRotationWorker()
}

// initializeDB initializes the database connection for the current month
func (al *AsyncLogger) initializeDB() error {
	al.dbMutex.Lock()
	defer al.dbMutex.Unlock()

	// Ensure logs directory exists
	if err := EnsureLogsDirectory(); err != nil {
		return err
	}

	// Get current month
	now := time.Now()
	currentMonth := fmt.Sprintf("%d-%02d", now.Year(), now.Month())

	// Close existing connection if month changed
	if al.db != nil && al.currentMonth != currentMonth {
		al.db.Close()
		al.db = nil
	}

	// Create new connection if needed
	if al.db == nil {
		logDbPath := utils.GetLogDbPath(currentMonth)
		db, err := initLogDB(logDbPath)
		if err != nil {
			return err
		}
		al.db = db
		al.currentMonth = currentMonth
	}

	return nil
}

// worker processes log entries asynchronously
func (al *AsyncLogger) worker() {
	defer al.wg.Done()

	batch := make([]LogEntry, 0, al.batchSize)
	ticker := time.NewTicker(al.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-al.ctx.Done():
			// Process remaining entries before shutting down
			if len(batch) > 0 {
				al.processBatch(batch)
			}
			return

		case entry := <-al.logChannel:
			batch = append(batch, entry)
			if len(batch) >= al.batchSize {
				al.processBatch(batch)
				batch = batch[:0] // Reset slice but keep capacity
			}

		case <-ticker.C:
			// Flush on timer
			if len(batch) > 0 {
				al.processBatch(batch)
				batch = batch[:0] // Reset slice but keep capacity
			}
		}
	}
}

// dbRotationWorker handles monthly database rotation
func (al *AsyncLogger) dbRotationWorker() {
	defer al.wg.Done()

	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	defer ticker.Stop()

	for {
		select {
		case <-al.ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			currentMonth := fmt.Sprintf("%d-%02d", now.Year(), now.Month())
			
			al.dbMutex.RLock()
			needsRotation := al.currentMonth != currentMonth
			al.dbMutex.RUnlock()

			if needsRotation {
				if err := al.initializeDB(); err != nil {
					log.Printf("Failed to rotate database: %v", err)
				}
			}
		}
	}
}

// processBatch writes a batch of log entries to the database
func (al *AsyncLogger) processBatch(batch []LogEntry) {
	if len(batch) == 0 {
		return
	}

	al.dbMutex.RLock()
	db := al.db
	al.dbMutex.RUnlock()

	if db == nil {
		log.Printf("Database not available, dropping %d log entries", len(batch))
		return
	}

	// Prepare batch insert statement
	insertSQL := `INSERT INTO access_logs (date, time, short_url, target_url, http_status, caller_ip, user_agent, referer, browser, browser_version, os, device_type, is_bot, ip_country, ip_city, ip_asn)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		log.Printf("Failed to prepare batch insert statement: %v", err)
		return
	}
	defer stmt.Close()

	// Execute batch insert using transaction for better performance
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback() // Will be no-op if transaction is committed

	txStmt := tx.Stmt(stmt)
	defer txStmt.Close()

	successCount := 0
	for _, entry := range batch {
		_, err := txStmt.Exec(entry.Date, entry.Time, entry.ShortURL, entry.TargetURL, 
			entry.HTTPStatus, entry.CallerIP, entry.UserAgent, entry.Referer, 
			entry.Browser, entry.BrowserVersion, entry.OS, entry.DeviceType, entry.IsBot,
			entry.IPCountry, entry.IPCity, entry.IPASN)
		if err != nil {
			log.Printf("Failed to insert log entry: %v", err)
			continue
		}
		successCount++
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit batch insert: %v", err)
		return
	}

	if successCount != len(batch) {
		log.Printf("Batch insert completed with %d/%d successful entries", successCount, len(batch))
	}
}

// LogAsync logs a link access asynchronously without blocking the HTTP request
func (al *AsyncLogger) LogAsync(r *http.Request, shortURL, targetURL string, httpStatus int) {
	// Create log entry quickly
	entry := al.createLogEntry(r, shortURL, targetURL, httpStatus)
	
	// Try to send to channel (non-blocking)
	select {
	case al.logChannel <- entry:
		// Successfully queued
	default:
		// Channel is full, drop the log entry (or implement a fallback strategy)
		log.Printf("Log channel full, dropping log entry for %s", shortURL)
	}
}

// createLogEntry creates a log entry from HTTP request data
func (al *AsyncLogger) createLogEntry(r *http.Request, shortURL, targetURL string, httpStatus int) LogEntry {
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

	// Parse User-Agent for detailed information (this is now async)
	uaInfo := ParseUserAgent(userAgent)

	// Parse IP address for geolocation information (this is now async)
	ipInfo := ParseIPAddress(clientIP)

	// Extract Referer
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "-"
	}

	return LogEntry{
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
		IPCountry:      ipInfo.Country,
		IPCity:         ipInfo.City,
		IPASN:          ipInfo.ASN,
	}
}

// Stop gracefully shuts down the async logger
func (al *AsyncLogger) Stop() {
	al.cancel()
	close(al.logChannel)
	al.wg.Wait()

	al.dbMutex.Lock()
	if al.db != nil {
		al.db.Close()
		al.db = nil
	}
	al.dbMutex.Unlock()
}

// LogLinkAccessAsync is the async version of LogLinkAccess that doesn't block
func LogLinkAccessAsync(r *http.Request, shortURL, targetURL string, httpStatus int) {
	logger := GetAsyncLogger()
	logger.LogAsync(r, shortURL, targetURL, httpStatus)
}