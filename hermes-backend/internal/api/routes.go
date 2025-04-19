package api

import (
	"net/http"

	"github.com/amaydixit11/hermes/hermes-backend/internal/api/handlers"
	"github.com/amaydixit11/hermes/hermes-backend/internal/api/middleware"
	"github.com/amaydixit11/hermes/hermes-backend/internal/config"
	"github.com/amaydixit11/hermes/hermes-backend/internal/service"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the HTTP routes for the API
func SetupRouter(cfg *config.Config, log *logger.Logger, serviceService *service.ServiceService) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.CORS())

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"version": "1.0.0",
			})
		})

		// Authentication routes
		// auth := v1.Group("/auth")
		// {
		// 	// These handlers will be implemented later
		// 	// auth.POST("/login", handlers.Login)
		// 	// auth.POST("/register", handlers.Register)
		// }

		// Protected routes
		// Will require implementing auth middleware
		protected := v1.Group("/")
		protected.Use(middleware.Auth(cfg))
		{
			// Service routes
			services := protected.Group("/services")
			{
				serviceHandler := handlers.NewServiceHandler(serviceService)
				services.POST("/", serviceHandler.RegisterService)
				services.GET("/", serviceHandler.ListServices)
				services.GET("/:id", serviceHandler.GetServiceByID)
				services.PUT("/:id", serviceHandler.UpdateService)
				services.DELETE("/:id", serviceHandler.DeleteService)
				// services.POST("/:id/health", serviceHandler.UpdateServiceHealth)
			}

			// // Gateway routes
			// gateway := protected.Group("/gateway")
			// {
			// 	// These handlers will be implemented later
			// 	// gateway.GET("/routes", handlers.ListRoutes)
			// 	// gateway.POST("/routes", handlers.CreateRoute)
			// }

			// // Metrics routes
			// metrics := protected.Group("/metrics")
			// {
			// 	// These handlers will be implemented later
			// 	// metrics.GET("/", handlers.GetMetrics)
			// 	// metrics.GET("/:service_id", handlers.GetServiceMetrics)
			// }
		}
	}

	return router
}
