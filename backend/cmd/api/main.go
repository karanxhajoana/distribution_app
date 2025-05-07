package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pack-sizer/internal/handlers"
	"pack-sizer/internal/services"
)

func main() {
	r := gin.Default()
	packSizesManager := services.NewPackSizesManager()
	packHandler := handlers.NewPackHandler(packSizesManager)

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Pack sizes endpoints
	r.GET("/api/pack-sizes", packHandler.GetPackSizes)
	r.PUT("/api/pack-sizes", packHandler.UpdatePackSize)
	r.POST("/api/pack-sizes", packHandler.AddPackSize)
	r.DELETE("/api/pack-sizes/:size", packHandler.RemovePackSize)

	// Calculate packs endpoint
	r.GET("/api/calculate", packHandler.CalculatePacks)

	r.Run(":8080")
} 