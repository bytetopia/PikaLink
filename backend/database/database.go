package database

import (
    "database/sql"
    "log"
    "os"
    "time"
    _ "modernc.org/sqlite"  // Replace the mattn/go-sqlite3 import
    "golang.org/x/crypto/bcrypt"
    "pikalink-backend/utils"
)

var DB *sql.DB

func InitDB() {
    // Get data path and create data directory if it doesn't exist
    dataPath := utils.GetDataPath()
    if err := os.MkdirAll(dataPath, 0755); err != nil {
        log.Fatal("Failed to create data directory:", err)
    }

    // Get database file path from centralized utils
    dbPath := utils.GetMainDbPath()
    log.Printf("Initializing database at: %s", dbPath)

    var err error
    DB, err = sql.Open("sqlite", dbPath)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    // Configure connection pool for optimal performance
    DB.SetMaxOpenConns(25)                     // Limit concurrent connections
    DB.SetMaxIdleConns(5)                      // Keep some connections alive
    DB.SetConnMaxLifetime(5 * time.Minute)     // Recycle old connections

    // Enable WAL mode and optimize SQLite settings for better concurrency
    if _, err := DB.Exec("PRAGMA journal_mode=WAL"); err != nil {
        log.Printf("Warning: Failed to enable WAL mode: %v", err)
    }
    if _, err := DB.Exec("PRAGMA synchronous=NORMAL"); err != nil {
        log.Printf("Warning: Failed to set synchronous mode: %v", err)
    }
    if _, err := DB.Exec("PRAGMA cache_size=1000"); err != nil {
        log.Printf("Warning: Failed to set cache size: %v", err)
    }
    if _, err := DB.Exec("PRAGMA temp_store=MEMORY"); err != nil {
        log.Printf("Warning: Failed to set temp store: %v", err)
    }

    createTables()
    log.Println("Database connected, optimized, and tables created")
}

func createTables() {
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`

    linkTable := `
    CREATE TABLE IF NOT EXISTS links (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        original_url TEXT NOT NULL,
        short_code TEXT UNIQUE NOT NULL,
        title TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        user_id INTEGER,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );`

    systemConfigTable := `
    CREATE TABLE IF NOT EXISTS system_config (
        config_key TEXT PRIMARY KEY,
        config_value TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

    if _, err := DB.Exec(userTable); err != nil {
        log.Fatal("Failed to create users table:", err)
    }

    if _, err := DB.Exec(linkTable); err != nil {
        log.Fatal("Failed to create links table:", err)
    }

    if _, err := DB.Exec(systemConfigTable); err != nil {
        log.Fatal("Failed to create system_config table:", err)
    }

    // Create default admin user (password: admin123)
    // Generate proper bcrypt hash for "admin123"
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal("Failed to hash password:", err)
    }
    
    DB.Exec("INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)", "admin", string(hashedPassword))
}