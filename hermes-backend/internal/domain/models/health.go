// models/health.go
package models

import (
	"time"
)

// HealthCheckType defines how a health check is performed
type HealthCheckType string

const (
	HealthCheckTypeActive  HealthCheckType = "ACTIVE"  // Hermes probes the service
	HealthCheckTypePassive HealthCheckType = "PASSIVE" // Service reports its own health
)

// HealthCheck represents a health check configuration for a service
type HealthCheck struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	ServiceID      string          `json:"service_id" gorm:"index;not null"`
	Name           string          `json:"name" gorm:"not null"`
	Type           HealthCheckType `json:"type" gorm:"not null"`
	Endpoint       string          `json:"endpoint"`
	Interval       int             `json:"interval" gorm:"default:60"`
	Timeout        int             `json:"timeout" gorm:"default:5"`
	Method         string          `json:"method" gorm:"default:GET"`
	ExpectedStatus int             `json:"expected_status"`
	ExpectedBody   string          `json:"expected_body"`
	Headers        string          `json:"headers" gorm:"type:jsonb"`
	Retries        int             `json:"retries" gorm:"default:1"`
	ThresholdCount int             `json:"threshold_count" gorm:"default:3"`
	TimeoutCount   int             `json:"timeout_count" gorm:"default:0"`
	Enabled        bool            `json:"enabled" gorm:"default:true"`
	CreatedAt      time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// HealthHistory stores historical health data for a service
type HealthHistory struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	ServiceID      string        `json:"service_id" gorm:"index;not null"`
	CheckID        uint          `json:"check_id" gorm:"index"`
	Status         ServiceStatus `json:"status" gorm:"not null"`
	Message        string        `json:"message"`
	ResponseTimeMs int           `json:"response_time_ms"`
	StatusCode     int           `json:"status_code"`
	Details        string        `json:"details" gorm:"type:jsonb"`
	Timestamp      time.Time     `json:"timestamp" gorm:"index;not null"`
}

// CustomHealthMetric represents a custom metric for service health
type CustomHealthMetric struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	ServiceID         string    `json:"service_id" gorm:"index;not null"`
	Name              string    `json:"name" gorm:"not null"`
	Description       string    `json:"description"`
	Value             float64   `json:"value"`
	Unit              string    `json:"unit"`
	WarningThreshold  float64   `json:"warning_threshold"`
	CriticalThreshold float64   `json:"critical_threshold"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// HealthThreshold defines when service status changes based on health checks
type HealthThreshold struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	ServiceID      string    `json:"service_id" gorm:"index;not null"`
	MetricName     string    `json:"metric_name"`
	WarningValue   float64   `json:"warning_value"`
	CriticalValue  float64   `json:"critical_value"`
	ComparisonType string    `json:"comparison_type" gorm:"default:'GREATER_THAN'"` // GREATER_THAN, LESS_THAN, EQUAL_TO
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Request/response models

// HealthCheckRequest represents the data needed to create a health check
type HealthCheckRequest struct {
	Name           string            `json:"name" binding:"required"`
	Type           HealthCheckType   `json:"type" binding:"required"`
	Endpoint       string            `json:"endpoint"`
	Interval       int               `json:"interval"`
	Timeout        int               `json:"timeout"`
	Method         string            `json:"method"`
	ExpectedStatus int               `json:"expected_status"`
	ExpectedBody   string            `json:"expected_body"`
	Headers        map[string]string `json:"headers"`
	Retries        int               `json:"retries"`
	ThresholdCount int               `json:"threshold_count"`
	Enabled        bool              `json:"enabled"`
}

// HealthUpdateRequest represents a health status update from a service
type HealthUpdateRequest struct {
	Status         ServiceStatus        `json:"status" binding:"required"`
	Message        string               `json:"message"`
	ResponseTimeMs int                  `json:"response_time_ms"`
	StatusCode     int                  `json:"status_code"`
	Details        map[string]string    `json:"details"`
	Metrics        []CustomMetricUpdate `json:"metrics"`
}

// CustomMetricUpdate represents an update to a custom health metric
type CustomMetricUpdate struct {
	Name  string  `json:"name" binding:"required"`
	Value float64 `json:"value" binding:"required"`
	Unit  string  `json:"unit"`
}

// HealthThresholdRequest represents the data needed to create a health threshold
type HealthThresholdRequest struct {
	MetricName     string  `json:"metric_name" binding:"required"`
	WarningValue   float64 `json:"warning_value"`
	CriticalValue  float64 `json:"critical_value"`
	ComparisonType string  `json:"comparison_type"`
}

// HealthHistoryQueryParams represents query parameters for listing health history
type HealthHistoryQueryParams struct {
	StartTime time.Time `form:"start_time"`
	EndTime   time.Time `form:"end_time"`
	Status    string    `form:"status"`
	Limit     int       `form:"limit,default=100"`
	Offset    int       `form:"offset,default=0"`
}
