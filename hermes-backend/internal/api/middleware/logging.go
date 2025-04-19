package middleware

import (
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // Process request

		duration := time.Since(start)

		log.Info("Incoming request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration.String(),
			"clientIP", c.ClientIP(),
		)
	}
}
