// service/health_service.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/repository"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
)

type HealthService struct {
	healthRepo  repository.HealthRepository
	serviceRepo repository.ServiceRepository
	log         *logger.Logger
	httpClient  *http.Client
}

func NewHealthService(
	healthRepo repository.HealthRepository,
	serviceRepo repository.ServiceRepository,
	log *logger.Logger,
) *HealthService {
	return &HealthService{
		healthRepo:  healthRepo,
		serviceRepo: serviceRepo,
		log:         log,
		httpClient:  &http.Client{},
	}
}

// Health check management
func (s *HealthService) CreateHealthCheck(ctx context.Context, serviceID string, req models.HealthCheckRequest) (*models.HealthCheck, error) {
	// Validate service exists
	_, err := s.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	// Create health check
	headersJSON, _ := json.Marshal(req.Headers)
	check := &models.HealthCheck{
		ServiceID:      serviceID,
		Name:           req.Name,
		Type:           req.Type,
		Endpoint:       req.Endpoint,
		Interval:       req.Interval,
		Timeout:        req.Timeout,
		Method:         req.Method,
		ExpectedStatus: req.ExpectedStatus,
		ExpectedBody:   req.ExpectedBody,
		Headers:        string(headersJSON),
		Retries:        req.Retries,
		ThresholdCount: req.ThresholdCount,
		Enabled:        req.Enabled,
	}

	err = s.healthRepo.CreateHealthCheck(ctx, check)
	if err != nil {
		return nil, fmt.Errorf("failed to create health check: %w", err)
	}

	s.log.Info("Created health check id=%d for service=%s", check.ID, serviceID)
	return check, nil
}

func (s *HealthService) GetHealthChecks(ctx context.Context, serviceID string) ([]*models.HealthCheck, error) {
	return s.healthRepo.GetHealthChecks(ctx, serviceID)
}
func (s *HealthService) GetHealthCheck(ctx context.Context, id uint) (*models.HealthCheck, error) {
	return s.healthRepo.GetHealthCheck(ctx, id)
}

func (s *HealthService) UpdateHealthCheck(ctx context.Context, id uint, req models.HealthCheckRequest) (*models.HealthCheck, error) {
	check, err := s.healthRepo.GetHealthCheck(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("health check not found: %w", err)
	}

	// Update fields
	if req.Name != "" {
		check.Name = req.Name
	}
	if req.Type != "" {
		check.Type = req.Type
	}
	if req.Endpoint != "" {
		check.Endpoint = req.Endpoint
	}
	if req.Interval > 0 {
		check.Interval = req.Interval
	}
	if req.Timeout > 0 {
		check.Timeout = req.Timeout
	}
	if req.Method != "" {
		check.Method = req.Method
	}
	if req.ExpectedStatus > 0 {
		check.ExpectedStatus = req.ExpectedStatus
	}
	if req.ExpectedBody != "" {
		check.ExpectedBody = req.ExpectedBody
	}
	if req.Headers != nil {
		headersJSON, _ := json.Marshal(req.Headers)
		check.Headers = string(headersJSON)
	}
	if req.Retries > 0 {
		check.Retries = req.Retries
	}
	if req.ThresholdCount > 0 {
		check.ThresholdCount = req.ThresholdCount
	}

	// Special case for boolean field to allow setting to false
	check.Enabled = req.Enabled

	err = s.healthRepo.UpdateHealthCheck(ctx, check)
	if err != nil {
		return nil, fmt.Errorf("failed to update health check: %w", err)
	}

	s.log.Info("Updated health check id=%d", id)
	return check, nil
}

func (s *HealthService) DeleteHealthCheck(ctx context.Context, id uint) error {
	return s.healthRepo.DeleteHealthCheck(ctx, id)
}

// Health reporting
func (s *HealthService) ReportServiceHealth(ctx context.Context, serviceID string, req models.HealthUpdateRequest) error {
	// Record health history
	history := &models.HealthHistory{
		ServiceID:      serviceID,
		Status:         req.Status,
		Message:        req.Message,
		ResponseTimeMs: req.ResponseTimeMs,
		StatusCode:     req.StatusCode,
		Timestamp:      time.Now(),
	}

	if req.Details != nil {
		detailsJSON, _ := json.Marshal(req.Details)
		history.Details = string(detailsJSON)
	}

	// Update service status
	err := s.healthRepo.UpdateHealthStatus(ctx, serviceID, req.Status, req.Message)
	if err != nil {
		return fmt.Errorf("failed to update service status: %w", err)
	}

	// Record health history
	err = s.healthRepo.RecordHealthHistory(ctx, history)
	if err != nil {
		return fmt.Errorf("failed to record health history: %w", err)
	}

	// Update custom metrics if provided
	if len(req.Metrics) > 0 {
		for _, metricUpdate := range req.Metrics {
			metric := &models.CustomHealthMetric{
				ServiceID: serviceID,
				Name:      metricUpdate.Name,
				Value:     metricUpdate.Value,
				Unit:      metricUpdate.Unit,
			}
			err = s.healthRepo.UpdateCustomMetric(ctx, metric)
			if err != nil {
				s.log.Error("Failed to update custom metric %s: %v", metricUpdate.Name, err)
			}
		}
	}

	s.log.Info("Updated health status for service=%s to %s", serviceID, req.Status)
	return nil
}

func (s *HealthService) GetHealthHistory(ctx context.Context, serviceID string, params models.HealthHistoryQueryParams) ([]*models.HealthHistory, int64, error) {
	return s.healthRepo.GetHealthHistory(ctx, serviceID, params)
}

// Custom metrics management
func (s *HealthService) GetCustomMetrics(ctx context.Context, serviceID string) ([]*models.CustomHealthMetric, error) {
	return s.healthRepo.GetCustomMetrics(ctx, serviceID)
}

func (s *HealthService) CreateOrUpdateMetric(ctx context.Context, serviceID string, metric models.CustomMetricUpdate) error {
	customMetric := &models.CustomHealthMetric{
		ServiceID: serviceID,
		Name:      metric.Name,
		Value:     metric.Value,
		Unit:      metric.Unit,
	}
	return s.healthRepo.UpdateCustomMetric(ctx, customMetric)
}

// Health threshold management
func (s *HealthService) CreateHealthThreshold(ctx context.Context, serviceID string, req models.HealthThresholdRequest) (*models.HealthThreshold, error) {
	threshold := &models.HealthThreshold{
		ServiceID:      serviceID,
		MetricName:     req.MetricName,
		WarningValue:   req.WarningValue,
		CriticalValue:  req.CriticalValue,
		ComparisonType: req.ComparisonType,
	}

	err := s.healthRepo.CreateHealthThreshold(ctx, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to create health threshold: %w", err)
	}

	return threshold, nil
}

func (s *HealthService) GetHealthThresholds(ctx context.Context, serviceID string) ([]*models.HealthThreshold, error) {
	return s.healthRepo.GetHealthThresholds(ctx, serviceID)
}

func (s *HealthService) UpdateHealthThreshold(ctx context.Context, id uint, req models.HealthThresholdRequest) (*models.HealthThreshold, error) {
	// Get the existing threshold
	thresholds, err := s.healthRepo.GetHealthThresholds(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve health thresholds: %w", err)
	}

	var threshold *models.HealthThreshold
	for _, t := range thresholds {
		if t.ID == id {
			threshold = t
			break
		}
	}

	if threshold == nil {
		return nil, fmt.Errorf("health threshold with ID %d not found", id)
	}

	// Update fields
	if req.MetricName != "" {
		threshold.MetricName = req.MetricName
	}
	threshold.WarningValue = req.WarningValue
	threshold.CriticalValue = req.CriticalValue
	if req.ComparisonType != "" {
		threshold.ComparisonType = req.ComparisonType
	}

	err = s.healthRepo.UpdateHealthThreshold(ctx, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to update health threshold: %w", err)
	}

	return threshold, nil
}

func (s *HealthService) DeleteHealthThreshold(ctx context.Context, id uint) error {
	return s.healthRepo.DeleteHealthThreshold(ctx, id)
}

// Active health checking
func (s *HealthService) RunActiveHealthCheck(ctx context.Context, check *models.HealthCheck) error {
	s.log.Debug("Running active health check for check_id=%d service_id=%s", check.ID, check.ServiceID)

	// Create HTTP request with timeout
	reqCtx, cancel := context.WithTimeout(ctx, time.Duration(check.Timeout)*time.Second)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(reqCtx, check.Method, check.Endpoint, nil)
	if err != nil {
		s.log.Error("Failed to create request: %v", err)
		return s.handleHealthCheckFailure(ctx, check)
	}

	// Add headers if specified
	if check.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(check.Headers), &headers); err == nil {
			for key, value := range headers {
				req.Header.Add(key, value)
			}
		}
	}

	// Execute request and measure time
	startTime := time.Now()
	resp, err := s.httpClient.Do(req)
	responseTime := time.Since(startTime).Milliseconds()

	// Handle request errors
	if err != nil {
		s.log.Error("Health check request failed: %v", err)
		return s.handleHealthCheckFailure(ctx, check)
	}
	defer resp.Body.Close()

	// Check if response meets expected status code
	if check.ExpectedStatus > 0 && resp.StatusCode != check.ExpectedStatus {
		s.log.Warn("Health check status code mismatch: expected=%d got=%d", check.ExpectedStatus, resp.StatusCode)
		return s.handleHealthCheckFailure(ctx, check)
	}

	// TODO: Check expected body if configured

	// Health check passed, reset failure count
	err = s.healthRepo.ResetHealthCheckFailures(ctx, check.ID)
	if err != nil {
		s.log.Error("Failed to reset health check failures: %v", err)
	}

	// Record successful health check in history
	history := &models.HealthHistory{
		ServiceID:      check.ServiceID,
		CheckID:        check.ID,
		Status:         models.ServiceStatusHealthy,
		ResponseTimeMs: int(responseTime),
		StatusCode:     resp.StatusCode,
		Timestamp:      time.Now(),
	}
	err = s.healthRepo.RecordHealthHistory(ctx, history)
	if err != nil {
		s.log.Error("Failed to record health history: %v", err)
	}

	// Update service status to healthy if needed
	service, err := s.serviceRepo.GetByID(ctx, check.ServiceID)
	if err != nil {
		s.log.Error("Failed to get service: %v", err)
		return nil
	}

	if service.Status != models.ServiceStatusHealthy {
		s.log.Info("Service %s recovered, updating status to HEALTHY", service.ID)
		err = s.healthRepo.UpdateHealthStatus(ctx, check.ServiceID, models.ServiceStatusHealthy, "Service recovered")
		if err != nil {
			s.log.Error("Failed to update service status: %v", err)
		}
	}

	return nil
}

func (s *HealthService) handleHealthCheckFailure(ctx context.Context, check *models.HealthCheck) error {
	// Increment failure count
	err := s.healthRepo.IncrementHealthCheckFailures(ctx, check.ID)
	if err != nil {
		s.log.Error("Failed to increment health check failures: %v", err)
		return err
	}

	// Get updated check to see if we've exceeded threshold
	updatedCheck, err := s.healthRepo.GetHealthCheck(ctx, check.ID)
	if err != nil {
		s.log.Error("Failed to get updated health check: %v", err)
		return err
	}

	// Record failed health check in history
	history := &models.HealthHistory{
		ServiceID: check.ServiceID,
		CheckID:   check.ID,
		Status:    models.ServiceStatusUnhealthy,
		Message:   "Health check failed",
		Timestamp: time.Now(),
	}
	err = s.healthRepo.RecordHealthHistory(ctx, history)
	if err != nil {
		s.log.Error("Failed to record health history: %v", err)
	}

	// If threshold exceeded, mark service as unhealthy
	if updatedCheck.TimeoutCount >= updatedCheck.ThresholdCount {
		s.log.Warn("Service %s health check failures exceeded threshold (%d/%d), marking as UNHEALTHY",
			check.ServiceID, updatedCheck.TimeoutCount, updatedCheck.ThresholdCount)

		err = s.healthRepo.UpdateHealthStatus(ctx, check.ServiceID, models.ServiceStatusUnhealthy,
			fmt.Sprintf("Health check '%s' failed %d times", check.Name, updatedCheck.TimeoutCount))
		if err != nil {
			s.log.Error("Failed to update service status: %v", err)
			return err
		}
	}

	return nil
}
