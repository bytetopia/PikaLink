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
    
    // Public routes
    r.POST("/api/login", handlers.Login)
    r.GET("/s/:code", handlers.RedirectLink)
    
    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware())
    {
        api.GET("/links", handlers.GetLinks)
        api.POST("/links", handlers.CreateLink)
        api.PUT("/links/:id", handlers.UpdateLink)
        api.DELETE("/links/:id", handlers.DeleteLink)
    }
    
    log.Println("Server starting on :8080")
    r.Run(":8080")
}