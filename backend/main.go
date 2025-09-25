package main

import (
    "log"
    "os"
    "strings"
    "pikalink-backend/database"
    "pikalink-backend/handlers"
    "pikalink-backend/middleware"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    // Initialize database
    database.InitDB()
    
    // Create Gin router
    r := gin.Default()
    
    // Get CORS origins from environment variable, default to localhost:3000 for development
    corsOrigins := os.Getenv("CORS_ORIGINS")
    if corsOrigins == "" {
        corsOrigins = "http://localhost:3000"
    }
    
    // Split the origins by comma for multiple origins
    allowedOrigins := strings.Split(corsOrigins, ",")
    for i, origin := range allowedOrigins {
        allowedOrigins[i] = strings.TrimSpace(origin)
    }
    
    // CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
    
    // API routes (unchanged)
    r.POST("/api/login", handlers.Login)
    
    // Protected API routes
    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware())
    {
        api.GET("/links", handlers.GetLinks)
        api.POST("/links", handlers.CreateLink)
        api.PUT("/links/:id", handlers.UpdateLink)
        api.DELETE("/links/:id", handlers.DeleteLink)
        api.POST("/change-password", handlers.ChangePassword)
        api.POST("/import", handlers.ImportLinks)
        api.GET("/export", handlers.ExportLinks)
    }
    
    // Serve admin frontend static files at /admin
    r.Static("/admin", "./frontend/")
    
    // Explicit route for root path to serve home page
    r.GET("/", func(c *gin.Context) {
        c.File("./frontend/home.html")
    })
    
    // Add specific route for short URL redirects (AFTER static routes)
    r.GET("/:code", func(c *gin.Context) {
        code := c.Param("code")
        // Special case: if code is "admin", redirect to admin page
        if code == "admin" {
            c.Redirect(302, "/admin/")
            return
        }
        handlers.RedirectLink(c)
    })
    
    // Handle frontend routes that don't match static files
    r.NoRoute(func(c *gin.Context) {
        path := c.Request.URL.Path
        // If it's an admin route that doesn't match a static file, serve index.html
        if strings.HasPrefix(path, "/admin") {
            c.File("./frontend/index.html")
            return
        }
        // Default 404
        c.JSON(404, gin.H{"error": "Not found"})
    })
    
    log.Println("Server starting on :8080")
    r.Run(":8080")
}