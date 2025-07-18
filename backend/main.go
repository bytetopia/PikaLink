package main

import (
    "log"
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
    
    // CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
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
    }
    
    // Serve admin frontend static files at /admin (MUST come before /:code route)
    r.Static("/admin", "./frontend/")
    
    // Add specific route for short URL redirects (AFTER static routes)
    r.GET("/:code", handlers.RedirectLink)
    
    // Handle frontend routes that don't match static files
    r.NoRoute(func(c *gin.Context) {
        path := c.Request.URL.Path
        
        // If it's an admin route that doesn't match a static file, serve index.html
        if len(path) > 6 && path[:6] == "/admin" {
            c.File("./frontend/index.html")
            return
        }
        
        // Default 404
        c.JSON(404, gin.H{"error": "Not found"})
    })
    
    log.Println("Server starting on :8080")
    r.Run(":8080")
}