package handlers

import (
    "crypto/rand"
    "database/sql"
    "encoding/base64"
    "encoding/csv"
    "fmt"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"
    "pikalink-backend/database"
    "pikalink-backend/models"
    "pikalink-backend/middleware"
    "pikalink-backend/logging"
    "pikalink-backend/utils"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

func generateShortCode() string {
    bytes := make([]byte, 6)
    rand.Read(bytes)
    return base64.URLEncoding.EncodeToString(bytes)[:8]
}

func isValidCustomCode(code string) bool {
    // Only allow alphanumeric characters, hyphens, and underscores
    // Length between 1 and 50 characters
    if len(code) < 1 || len(code) > 50 {
        return false
    }
    
    matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", code)
    return matched
}

func isCustomCodeAvailable(code string) bool {
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM links WHERE short_code = ?", code).Scan(&count)
    return err == nil && count == 0
}

func Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    err := database.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", req.Username).
        Scan(&user.ID, &user.Username, &user.Password)
    
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := middleware.GenerateJWT(user.ID, user.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Check if user is using default password (admin123)
    isDefaultPassword := req.Password == "admin123"

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user":  gin.H{"id": user.ID, "username": user.Username, "is_default_password": isDefaultPassword},
    })
}

func CreateLink(c *gin.Context) {
    var req models.CreateLinkRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.MustGet("user_id").(int)
    var shortCode string

    if req.UseCustomCode {
        // Validate custom code
        customCode := strings.TrimSpace(req.CustomCode)
        if customCode == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Custom code cannot be empty"})
            return
        }
        
        if !isValidCustomCode(customCode) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Custom code must be 1-50 characters long and contain only letters, numbers, hyphens, and underscores"})
            return
        }
        
        if !isCustomCodeAvailable(customCode) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Custom code is already taken"})
            return
        }
        
        shortCode = customCode
    } else {
        // Generate random code
        shortCode = generateShortCode()
        
        // Ensure generated code is unique
        for !isCustomCodeAvailable(shortCode) {
            shortCode = generateShortCode()
        }
    }

    query := `INSERT INTO links (original_url, short_code, title, user_id, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?)`

    now := time.Now()
    result, err := database.DB.Exec(query, req.OriginalURL, shortCode, req.Title, userID, now, now)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
        return
    }

    id, _ := result.LastInsertId()
    link := models.Link{
        ID:          int(id),
        OriginalURL: req.OriginalURL,
        ShortCode:   shortCode,
        Title:       req.Title,
        UserID:      userID,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    c.JSON(http.StatusCreated, link)
}

func GetLinks(c *gin.Context) {
    userID := c.MustGet("user_id").(int)
    
    rows, err := database.DB.Query(`
        SELECT id, original_url, short_code, title, created_at, updated_at 
        FROM links WHERE user_id = ? ORDER BY created_at DESC
    `, userID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
        return
    }
    defer rows.Close()

    var links []models.Link
    for rows.Next() {
        var link models.Link
        err := rows.Scan(&link.ID, &link.OriginalURL, &link.ShortCode, &link.Title, 
                        &link.CreatedAt, &link.UpdatedAt)
        if err != nil {
            continue
        }
        link.UserID = userID
        links = append(links, link)
    }

    c.JSON(http.StatusOK, links)
}

func GetLink(c *gin.Context) {
    id := c.Param("id")
    linkID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
        return
    }

    userID := c.MustGet("user_id").(int)
    
    var link models.Link
    err = database.DB.QueryRow(`
        SELECT id, original_url, short_code, title, created_at, updated_at 
        FROM links WHERE id = ? AND user_id = ?
    `, linkID, userID).Scan(&link.ID, &link.OriginalURL, &link.ShortCode, &link.Title, 
                            &link.CreatedAt, &link.UpdatedAt)
    
    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
        return
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch link"})
        return
    }
    
    link.UserID = userID
    c.JSON(http.StatusOK, link)
}

func UpdateLink(c *gin.Context) {
    id := c.Param("id")
    linkID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
        return
    }

    var req models.UpdateLinkRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.MustGet("user_id").(int)
    
    // Get the current link to check if it exists and belongs to the user
    var currentLink models.Link
    err = database.DB.QueryRow("SELECT id, short_code FROM links WHERE id = ? AND user_id = ?", linkID, userID).
        Scan(&currentLink.ID, &currentLink.ShortCode)
    
    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
        return
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // Handle short code update
    var shortCodeToUpdate string
    if req.ShortCode != "" && req.ShortCode != currentLink.ShortCode {
        // Validate the new short code
        if !isValidCustomCode(req.ShortCode) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Short code must be 1-50 characters long and contain only letters, numbers, hyphens, and underscores"})
            return
        }
        
        // Check if the new short code is available
        if !isCustomCodeAvailable(req.ShortCode) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is already taken"})
            return
        }
        
        shortCodeToUpdate = req.ShortCode
    } else {
        shortCodeToUpdate = currentLink.ShortCode
    }
    
    query := `UPDATE links SET original_url = ?, title = ?, short_code = ?, updated_at = ? 
              WHERE id = ? AND user_id = ?`
    
    result, err := database.DB.Exec(query, req.OriginalURL, req.Title, shortCodeToUpdate, time.Now(), linkID, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update link"})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Link updated successfully"})
}

func DeleteLink(c *gin.Context) {
    id := c.Param("id")
    linkID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
        return
    }

    userID := c.MustGet("user_id").(int)
    
    result, err := database.DB.Exec("DELETE FROM links WHERE id = ? AND user_id = ?", linkID, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete link"})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}

func RedirectLink(c *gin.Context) {
    var shortCode string
    
    // Try to get from route parameter first (if called from a proper route)
    shortCode = c.Param("code")
    
    // If no route parameter, extract from URL path (when called from NoRoute)
    if shortCode == "" {
        path := c.Request.URL.Path
        if len(path) > 1 {
            shortCode = path[1:] // Remove leading slash
        }
    }

    if shortCode == "" {
        // Log the failed access attempt (async)
        logging.LogLinkAccessAsync(c.Request, shortCode, "", http.StatusNotFound)
        Serve404Page(c)
        return
    }

    var originalURL string
    var linkID int
    err := database.DB.QueryRow("SELECT id, original_url FROM links WHERE short_code = ?", shortCode).
        Scan(&linkID, &originalURL)

    if err == sql.ErrNoRows {
        // Log the failed access attempt (async)
        logging.LogLinkAccessAsync(c.Request, shortCode, "", http.StatusNotFound)
        Serve404Page(c)
        return
    }

    if err != nil {
        // Log the error (async)
        logging.LogLinkAccessAsync(c.Request, shortCode, "", http.StatusInternalServerError)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    
    // Merge the original URL with any additional path and query parameters from the short URL
    mergedURL, err := utils.MergeURLWithShortURL(originalURL, c.Request.URL.Path, c.Request.URL.RawQuery)
    if err != nil {
        // Log the error (async)
        logging.LogLinkAccessAsync(c.Request, shortCode, originalURL, http.StatusInternalServerError)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "URL merge error"})
        return
    }
    
    // Log the successful redirect (async) - this happens AFTER the redirect
    // so it doesn't block the user's redirect
    logging.LogLinkAccessAsync(c.Request, shortCode, mergedURL, http.StatusFound)
    
    c.Redirect(http.StatusFound, mergedURL)
}

func ChangePassword(c *gin.Context) {
    var req models.ChangePasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.MustGet("user_id").(int)

    // Validate new password length
    if len(req.NewPassword) < 6 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be at least 6 characters long"})
        return
    }

    // Get current user password from database
    var currentHashedPassword string
    err := database.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).
        Scan(&currentHashedPassword)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
        return
    }

    // Verify current password
    err = bcrypt.CompareHashAndPassword([]byte(currentHashedPassword), []byte(req.CurrentPassword))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
        return
    }

    // Hash new password
    newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Update password in database
    _, err = database.DB.Exec("UPDATE users SET password = ? WHERE id = ?", string(newHashedPassword), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
        return
    }

    // Check if the new password is still the default password
    isStillDefault := req.NewPassword == "admin123"

    c.JSON(http.StatusOK, gin.H{
        "message": "Password changed successfully",
        "is_default_password": isStillDefault,
    })
}

// ImportResult represents the result of a CSV import operation
type ImportResult struct {
    SuccessCount int      `json:"success_count"`
    FailedCount  int      `json:"failed_count"`
    FailedItems  []string `json:"failed_items"`
    TotalCount   int      `json:"total_count"`
}

// ImportLinks handles CSV import of links
func ImportLinks(c *gin.Context) {
    userID := c.MustGet("user_id").(int)
    
    // Get the uploaded file
    file, _, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }
    defer file.Close()

    // Parse CSV
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV format"})
        return
    }

    if len(records) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Empty CSV file"})
        return
    }

    // Validate header
    if len(records[0]) < 3 || records[0][0] != "original_url" || records[0][1] != "short_code" || records[0][2] != "title" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV header. Expected: original_url,short_code,title"})
        return
    }

    result := ImportResult{
        SuccessCount: 0,
        FailedCount:  0,
        FailedItems:  []string{},
        TotalCount:   len(records) - 1, // Exclude header
    }

    // Process each record (skip header)
    for i, record := range records[1:] {
        if len(record) < 3 {
            result.FailedCount++
            result.FailedItems = append(result.FailedItems, fmt.Sprintf("Row %d: insufficient columns", i+2))
            continue
        }

        originalURL := strings.TrimSpace(record[0])
        shortCode := strings.TrimSpace(record[1])
        title := strings.TrimSpace(record[2])

        // Validate required fields
        if originalURL == "" {
            result.FailedCount++
            result.FailedItems = append(result.FailedItems, fmt.Sprintf("Row %d: empty original_url", i+2))
            continue
        }

        if shortCode == "" {
            result.FailedCount++
            result.FailedItems = append(result.FailedItems, fmt.Sprintf("Row %d: empty short_code", i+2))
            continue
        }

        // Validate short code format
        if !isValidCustomCode(shortCode) {
            result.FailedCount++
            result.FailedItems = append(result.FailedItems, fmt.Sprintf("Row %d: invalid short_code '%s'", i+2, shortCode))
            continue
        }

        // Check if short code already exists
        if !isCustomCodeAvailable(shortCode) {
            result.FailedCount++
            result.FailedItems = append(result.FailedItems, fmt.Sprintf("Row %d: short_code '%s' already exists", i+2, shortCode))
            continue
        }

        // Insert the link
        query := `INSERT INTO links (original_url, short_code, title, user_id, created_at, updated_at) 
                  VALUES (?, ?, ?, ?, ?, ?)`
        
        now := time.Now()
        _, err := database.DB.Exec(query, originalURL, shortCode, title, userID, now, now)
        if err != nil {
            result.FailedCount++
            result.FailedItems = append(result.FailedItems, fmt.Sprintf("Row %d: database error for short_code '%s'", i+2, shortCode))
            continue
        }

        result.SuccessCount++
    }

    c.JSON(http.StatusOK, result)
}

// ExportLinks handles CSV export of all links for the current user
func ExportLinks(c *gin.Context) {
    userID := c.MustGet("user_id").(int)
    
    // Query all links for the user
    rows, err := database.DB.Query(`
        SELECT original_url, short_code, title 
        FROM links WHERE user_id = ? ORDER BY created_at DESC
    `, userID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
        return
    }
    defer rows.Close()

    // Set headers for CSV download
    c.Header("Content-Type", "text/csv")
    c.Header("Content-Disposition", "attachment; filename=pikalink_export.csv")

    // Create CSV writer
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()

    // Write header
    writer.Write([]string{"original_url", "short_code", "title"})

    // Write data
    for rows.Next() {
        var originalURL, shortCode, title string
        if err := rows.Scan(&originalURL, &shortCode, &title); err != nil {
            continue
        }
        
        writer.Write([]string{originalURL, shortCode, title})
    }

    if err = rows.Err(); err != nil {
        return
    }
}