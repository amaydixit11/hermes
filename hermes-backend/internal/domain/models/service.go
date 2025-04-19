package models

import (
	"time"

	"github.com/lib/pq"
)

// Service represents a microservice in the system
type Service struct {
	ID           string            `json:"id" gorm:"primaryKey"`
	Name         string            `json:"name" gorm:"uniqueIndex;not null"`
	Description  string            `json:"description"`
	Status       ServiceStatus     `json:"status" gorm:"not null;default:'UNKNOWN'"`
	Type         string            `json:"type"`
	Endpoint     string            `json:"endpoint" gorm:"not null"`
	Metadata     map[string]string `json:"metadata" gorm:"serializer:json"`
	Tags         pq.StringArray    `gorm:"type:text[]"`
	CreatedAt    time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	LastSeen     time.Time         `json:"last_seen"`
	RegisteredBy string            `json:"registered_by,omitempty"`
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
	Name         string            `json:"name" binding:"required"`
	Description  string            `json:"description"`
	Type         string            `json:"type"`
	Endpoint     string            `json:"endpoint" binding:"required"`
	Metadata     map[string]string `json:"metadata"`
	Tags         []string          `json:"tags"`
	RegisteredBy string            `json:"registered_by"`
}

type BulkServiceRegistration struct {
	Services       []ServiceRegistration `json:"services" binding:"required"`
	CommonMetadata map[string]string     `json:"common_metadata"`
	CommonTags     []string              `json:"common_tags"`
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

// ServiceVersion tracks different versions of a service
type ServiceVersion struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ServiceID   string    `json:"service_id" gorm:"index"`
	Version     string    `json:"version" gorm:"not null"`
	IsActive    bool      `json:"is_active" gorm:"default:false"`
	Endpoint    string    `json:"endpoint" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ServiceVersionRequest represents the data needed to register a new service version
type ServiceVersionRequest struct {
	Version     string `json:"version" binding:"required"`
	IsActive    bool   `json:"is_active"`
	Endpoint    string `json:"endpoint" binding:"required"`
	Description string `json:"description"`
}

// ServiceDependency tracks dependencies between services
type ServiceDependency struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	ServiceID      string    `json:"service_id" gorm:"index;not null"`
	DependencyID   string    `json:"dependency_id" gorm:"index;not null"`
	DependencyType string    `json:"dependency_type" gorm:"not null"` // e.g., "REQUIRED", "OPTIONAL"
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// ServiceDependencyRequest represents the data needed to register a new service dependency
type ServiceDependencyRequest struct {
	DependencyID   string `json:"dependency_id" binding:"required"`
	DependencyType string `json:"dependency_type" binding:"required"`
	Description    string `json:"description"`
}

// AdvancedDiscoveryParams represents extended query parameters for service discovery
type AdvancedDiscoveryParams struct {
	ServiceQueryParams

	HealthStatus []string `form:"health_status"`

	MetadataKey   string `form:"metadata_key"`
	MetadataValue string `form:"metadata_value"`

	IncludeVersions   bool   `form:"include_versions"`
	VersionFilter     string `form:"version"`
	ActiveVersionOnly bool   `form:"active_version_only"`

	DependencyID string `form:"dependency_id"`
	DependencyOf string `form:"dependency_of"`

	LastSeenSince time.Time `form:"last_seen_since"`
}
