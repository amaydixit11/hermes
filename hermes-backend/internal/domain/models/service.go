package models

import (
	"time"

	"github.com/lib/pq"
)

// Service represents a microservice in the system
type Service struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name" gorm:"uniqueIndex;not null"`
	Description string            `json:"description"`
	Status      ServiceStatus     `json:"status" gorm:"not null;default:'UNKNOWN'"`
	Type        string            `json:"type"`
	Endpoint    string            `json:"endpoint" gorm:"not null"`
	Metadata    map[string]string `json:"metadata" gorm:"serializer:json"`
	Tags        pq.StringArray    `gorm:"type:text[]"`
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	LastSeen    time.Time         `json:"last_seen"`
}

// ServiceStatus represents the health status of a service
type ServiceStatus string

// Service status constants
const (
	ServiceStatusHealthy   ServiceStatus = "HEALTHY"
	ServiceStatusUnhealthy ServiceStatus = "UNHEALTHY"
	ServiceStatusUnknown   ServiceStatus = "UNKNOWN"
	ServiceStatusWarning   ServiceStatus = "WARNING"
)

// ServiceRegistration represents the data needed to register a new service
type ServiceRegistration struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Endpoint    string            `json:"endpoint" binding:"required"`
	Metadata    map[string]string `json:"metadata"`
	Tags        []string          `json:"tags"`
}

// ServiceHealthUpdate represents a health status update from a service
type ServiceHealthUpdate struct {
	Status  ServiceStatus     `json:"status" binding:"required"`
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

// ServiceUpdateRequest represents the data that can be updated for a service
type ServiceUpdateRequest struct {
	Name        *string           `json:"name"`
	Description *string           `json:"description"`
	Status      *ServiceStatus    `json:"status"`
	Type        *string           `json:"type"`
	Endpoint    *string           `json:"endpoint"`
	Metadata    map[string]string `json:"metadata"`
	Tags        []string          `json:"tags"`
}

// ServiceQueryParams represents query parameters for listing services
type ServiceQueryParams struct {
	Status string   `form:"status"`
	Type   string   `form:"type"`
	Tags   []string `form:"tags"`
	Search string   `form:"search"`
	Limit  int      `form:"limit,default=20"`
	Offset int      `form:"offset,default=0"`
}
