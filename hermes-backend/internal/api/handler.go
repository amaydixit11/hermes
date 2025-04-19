// internal/api/handler.go
package api

import (
	"github.com/gin-gonic/gin"
)

// Register routes for the API
func RegisterRoutes(router *gin.Engine) {
	router.GET("/api/status", getStatus)
	// Add more routes here as needed
}

// Example of a basic route handler
func getStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "API is running",
	})
}
