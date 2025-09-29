package database

import (
    "database/sql"
    "log"
    "os"
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

    createTables()
    log.Println("Database connected and tables created")
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

    if _, err := DB.Exec(userTable); err != nil {
        log.Fatal("Failed to create users table:", err)
    }

    if _, err := DB.Exec(linkTable); err != nil {
        log.Fatal("Failed to create links table:", err)
    }

    // Create default admin user (password: admin123)
    // Generate proper bcrypt hash for "admin123"
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal("Failed to hash password:", err)
    }
    
    DB.Exec("INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)", "admin", string(hashedPassword))
}