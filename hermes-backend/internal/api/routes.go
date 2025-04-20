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
func SetupRouter(cfg *config.Config, log *logger.Logger, serviceService *service.ServiceService, healthService *service.HealthService) *gin.Engine {
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
		// These handlers will be implemented later
		// auth.POST("/login", handlers.Login)
		// auth.POST("/register", handlers.Register)
		// }

		// Protected routes
		protected := v1.Group("/")
		// TODO: Add Auth
		// protected.Use(middleware.Auth(cfg))
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
				services.POST("/bulk", serviceHandler.BulkRegisterService)
				// Add this to your existing routes setup

				// Service Discovery routes
				discoveryHandler := handlers.NewServiceDiscoveryHandler(serviceService)
				services.GET("/discovery", discoveryHandler.AdvancedSearch)

				// Service Version routes
				versionHandler := handlers.NewServiceVersionHandler(serviceService)
				services.POST("/:id/versions", versionHandler.AddServiceVersion)
				services.GET("/:id/versions", versionHandler.GetServiceVersions)
				services.PUT("/:id/versions/:version/activate", versionHandler.ActivateServiceVersion)

				// Service Dependency routes
				dependencyHandler := handlers.NewServiceDependencyHandler(serviceService)
				services.POST("/:id/dependencies", dependencyHandler.AddServiceDependency)
				services.GET("/:id/dependencies", dependencyHandler.GetServiceDependencies)
				services.GET("/:id/dependents", dependencyHandler.GetServiceDependents)
				services.DELETE("/:id/dependencies/:dependency_id", dependencyHandler.RemoveServiceDependency) // api/router.go (add to your existing routes)

				// Create health handler
				healthHandler := handlers.NewHealthHandler(healthService)

				// Health check routes
				services.POST("/:id/health", healthHandler.ReportServiceHealth)
				services.GET("/:id/health/history", healthHandler.GetHealthHistory)

				// Health checks configuration routes
				services.POST("/:id/health-checks", healthHandler.CreateHealthCheck)
				services.GET("/:id/health-checks", healthHandler.GetHealthChecks)
				services.GET("/:id/health-checks/:check_id", healthHandler.GetHealthCheck)
				services.PUT("/:id/health-checks/:check_id", healthHandler.UpdateHealthCheck)
				services.DELETE("/:id/health-checks/:check_id", healthHandler.DeleteHealthCheck)

				// Custom metrics routes
				services.GET("/:id/metrics", healthHandler.GetCustomMetrics)
				services.POST("/:id/metrics", healthHandler.CreateOrUpdateCustomMetric)

				// Health thresholds routes
				services.POST("/:id/thresholds", healthHandler.CreateHealthThreshold)
				services.GET("/:id/thresholds", healthHandler.GetHealthThresholds)
				services.PUT("/:id/thresholds/:threshold_id", healthHandler.UpdateHealthThreshold)
				services.DELETE("/:id/thresholds/:threshold_id", healthHandler.DeleteHealthThreshold)

			}

			// // Gateway routes
			// gateway := protected.Group("/gateway")
			// {
			// These handlers will be implemented later
			// gateway.GET("/routes", handlers.ListRoutes)
			// gateway.POST("/routes", handlers.CreateRoute)
			// }

			// // Metrics routes
			// metrics := protected.Group("/metrics")
			// {
			// These handlers will be implemented later
			// metrics.GET("/", handlers.GetMetrics)
			// metrics.GET("/:service_id", handlers.GetServiceMetrics)
			// }
		}
	}

	return router
}
