// cmd/main.go
package main

import (
	"log"

	"github.com/amaydixit11/hermes/hermes-backend/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Define API routes
	api.RegisterRoutes(router)

	// Start the server
	if err := router.Run(":4213"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
