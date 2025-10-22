package database

import (
	"database/sql"
	"sync"
	"time"

	"pikalink-backend/utils"
	_ "modernc.org/sqlite"
)

// AnalysisDBPool manages connections to monthly log databases for analysis
type AnalysisDBPool struct {
	pools map[string]*sql.DB
	mutex sync.RWMutex
}

var (
	analysisPool *AnalysisDBPool
	poolOnce     sync.Once
)

// GetAnalysisPool returns the singleton analysis database pool
func GetAnalysisPool() *AnalysisDBPool {
	poolOnce.Do(func() {
		analysisPool = &AnalysisDBPool{
			pools: make(map[string]*sql.DB),
		}
	})
	return analysisPool
}

// GetConnection returns a connection to the specified month's log database
func (p *AnalysisDBPool) GetConnection(month string) (*sql.DB, error) {
	// Try to get existing connection
	p.mutex.RLock()
	if db, exists := p.pools[month]; exists {
		p.mutex.RUnlock()
		// Test if connection is still alive
		if err := db.Ping(); err == nil {
			return db, nil
		}
		// Connection is dead, remove it and create a new one
		p.mutex.RUnlock()
		p.mutex.Lock()
		delete(p.pools, month)
		if db != nil {
			db.Close()
		}
		p.mutex.Unlock()
	} else {
		p.mutex.RUnlock()
	}

	// Create new connection
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Double-check after acquiring write lock
	if db, exists := p.pools[month]; exists {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		// Clean up dead connection
		delete(p.pools, month)
		if db != nil {
			db.Close()
		}
	}

	dbPath := utils.GetLogDbPath(month)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Configure for read-only analysis workload
	db.SetMaxOpenConns(5)                    // Limit concurrent connections per month
	db.SetMaxIdleConns(2)                    // Keep a few connections alive
	db.SetConnMaxLifetime(10 * time.Minute)  // Recycle connections every 10 minutes

	// Optimize for read-only operations
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		// WAL mode might not be available, continue without it
	}
	if _, err := db.Exec("PRAGMA query_only=ON"); err != nil {
		// Read-only mode might not be available in all SQLite versions
	}
	if _, err := db.Exec("PRAGMA cache_size=2000"); err != nil {
		// Larger cache for analysis queries
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	p.pools[month] = db
	return db, nil
}

// CloseAll closes all connections in the pool
func (p *AnalysisDBPool) CloseAll() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for month, db := range p.pools {
		if db != nil {
			db.Close()
		}
		delete(p.pools, month)
	}
}

// CloseMonth closes the connection for a specific month
func (p *AnalysisDBPool) CloseMonth(month string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if db, exists := p.pools[month]; exists {
		if db != nil {
			db.Close()
		}
		delete(p.pools, month)
	}
}

// GetStats returns statistics about the connection pool
func (p *AnalysisDBPool) GetStats() map[string]interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["total_connections"] = len(p.pools)
	
	months := make([]string, 0, len(p.pools))
	for month := range p.pools {
		months = append(months, month)
	}
	stats["active_months"] = months

	return stats
}