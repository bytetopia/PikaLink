package handlers

import (
    "database/sql"
    "net/http"
    "time"
    "pikalink-backend/database"
    "pikalink-backend/models"
    "pikalink-backend/defaults"
    "github.com/gin-gonic/gin"
)

// GetConfig retrieves a configuration value by key
func GetConfig(c *gin.Context) {
    key := c.Param("key")
    
    var config models.SystemConfig
    err := database.DB.QueryRow(
        "SELECT config_key, config_value, created_at, updated_at FROM system_config WHERE config_key = ?", 
        key,
    ).Scan(&config.ConfigKey, &config.ConfigValue, &config.CreatedAt, &config.UpdatedAt)
    
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration"})
        return
    }
    
    c.JSON(http.StatusOK, config)
}

// UpdateConfig updates a configuration value by key
func UpdateConfig(c *gin.Context) {
    key := c.Param("key")
    
    var request models.UpdateConfigRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    // Check if config exists
    var exists bool
    err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM system_config WHERE config_key = ?)", key).Scan(&exists)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check configuration"})
        return
    }
    
    if exists {
        // Update existing configuration
        _, err = database.DB.Exec(
            "UPDATE system_config SET config_value = ?, updated_at = ? WHERE config_key = ?",
            request.ConfigValue, time.Now(), key,
        )
    } else {
        // Create new configuration
        _, err = database.DB.Exec(
            "INSERT INTO system_config (config_key, config_value, created_at, updated_at) VALUES (?, ?, ?, ?)",
            key, request.ConfigValue, time.Now(), time.Now(),
        )
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

// GetPageContent retrieves HTML content from database or returns default
func GetPageContent(configKey string) (string, error) {
    var content string
    err := database.DB.QueryRow(
        "SELECT config_value FROM system_config WHERE config_key = ?", 
        configKey,
    ).Scan(&content)
    
    if err != nil {
        if err == sql.ErrNoRows {
            // Return default content based on key
            switch configKey {
            case models.ConfigKeyHomePage:
                return defaults.GetDefaultHomePage(), nil
            case models.ConfigKey404Page:
                return defaults.GetDefault404Page(), nil
            default:
                return "", err
            }
        }
        return "", err
    }
    
    return content, nil
}

// ServeHomePage serves the home page with content from database
func ServeHomePage(c *gin.Context) {
    content, err := GetPageContent(models.ConfigKeyHomePage)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load home page"})
        return
    }
    
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

// Serve404Page serves the 404 page with content from database
func Serve404Page(c *gin.Context) {
    content, err := GetPageContent(models.ConfigKey404Page)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load 404 page"})
        return
    }
    
    c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(content))
}

