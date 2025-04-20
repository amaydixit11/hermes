// handlers/health_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthService *service.HealthService
}

func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}
func (h *HealthHandler) CreateHealthCheck(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service ID is required"})
		return
	}
	// handlers/health_handler.go (continued)
	var req models.HealthCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	check, err := h.healthService.CreateHealthCheck(c, serviceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create health check: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, check)
}

func (h *HealthHandler) GetHealthChecks(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service ID is required"})
		return
	}

	checks, err := h.healthService.GetHealthChecks(c, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get health checks: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, checks)
}

func (h *HealthHandler) GetHealthCheck(c *gin.Context) {
	idStr := c.Param("check_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check ID is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check ID"})
		return
	}

	check, err := h.healthService.GetHealthCheck(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "health check not found"})
		return
	}

	c.JSON(http.StatusOK, check)
}

func (h *HealthHandler) UpdateHealthCheck(c *gin.Context) {
	idStr := c.Param("check_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check ID is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check ID"})
		return
	}

	var req models.HealthCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	check, err := h.healthService.UpdateHealthCheck(c, uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update health check: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, check)
}

func (h *HealthHandler) DeleteHealthCheck(c *gin.Context) {
	idStr := c.Param("check_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check ID is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check ID"})
		return
	}

	err = h.healthService.DeleteHealthCheck(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete health check: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "health check deleted successfully"})
}

func (h *HealthHandler) ReportServiceHealth(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service ID is required"})
		return
	}

	var req models.HealthUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	err := h.healthService.ReportServiceHealth(c, serviceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to report service health: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "health status updated successfully"})
}

func (h *HealthHandler) GetHealthHistory(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service ID is required"})
		return
	}

	var params models.HealthHistoryQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters: " + err.Error()})
		return
	}

	history, total, err := h.healthService.GetHealthHistory(c, serviceID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get health history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": history,
		"total": total,
	})
}

func (h *HealthHandler) GetCustomMetrics(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service ID is required"})
		return
	}

	metrics, err := h.healthService.GetCustomMetrics(c, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get custom metrics: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *HealthHandler) CreateOrUpdateCustomMetric(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service ID is required"})
		return
	}

	var req models.CustomMetricUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	err := h.healthService.CreateOrUpdateMetric(c, serviceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create/update metric: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "metric updated successfully"})
}

func (h *HealthHandler) CreateHealthThreshold(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service ID is required"})
		return
	}

	var req models.HealthThresholdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	threshold, err := h.healthService.CreateHealthThreshold(c, serviceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create health threshold: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, threshold)
}

func (h *HealthHandler) GetHealthThresholds(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service ID is required"})
		return
	}

	thresholds, err := h.healthService.GetHealthThresholds(c, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get health thresholds: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, thresholds)
}

func (h *HealthHandler) UpdateHealthThreshold(c *gin.Context) {
	idStr := c.Param("threshold_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "threshold ID is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid threshold ID"})
		return
	}

	var req models.HealthThresholdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	threshold, err := h.healthService.UpdateHealthThreshold(c, uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update health threshold: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, threshold)
}

func (h *HealthHandler) DeleteHealthThreshold(c *gin.Context) {
	idStr := c.Param("threshold_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "threshold ID is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid threshold ID"})
		return
	}

	err = h.healthService.DeleteHealthThreshold(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete health threshold: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "health threshold deleted successfully"})
}
