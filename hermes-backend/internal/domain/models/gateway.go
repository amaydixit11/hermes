package models

import (
	"time"
)

// Route represents an API gateway route configuration
type Route struct {
	ID             string            `json:"id" db:"id"`
	Path           string            `json:"path" db:"path"`
	Description    string            `json:"description" db:"description"`
	ServiceID      string            `json:"service_id" db:"service_id"`
	LoadBalancerID string            `json:"load_balancer_id" db:"load_balancer_id"`
	Targets        []string          `json:"targets" db:"targets"`
	Active         bool              `json:"active" db:"active"`
	Headers        map[string]string `json:"headers" db:"headers"`
	RateLimit      *RateLimit        `json:"rate_limit" db:"rate_limit"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
}

// RateLimit defines rate limiting configuration for a route
type RateLimit struct {
	Limit  int           `json:"limit"`  // Number of requests
	Window time.Duration `json:"window"` // Time window for the limit
	PerIP  bool          `json:"per_ip"` // Whether to apply per IP address
}

// RouteCreationRequest represents a request to create a new route
type RouteCreationRequest struct {
	Path           string            `json:"path" binding:"required"`
	Description    string            `json:"description"`
	ServiceID      string            `json:"service_id" binding:"required"`
	LoadBalancerID string            `json:"load_balancer_id"`
	Targets        []string          `json:"targets" binding:"required"`
	Headers        map[string]string `json:"headers"`
	RateLimit      *RateLimit        `json:"rate_limit"`
}

// RouteUpdateRequest represents a request to update an existing route
type RouteUpdateRequest struct {
	Path           *string           `json:"path"`
	Description    *string           `json:"description"`
	ServiceID      *string           `json:"service_id"`
	LoadBalancerID *string           `json:"load_balancer_id"`
	Targets        []string          `json:"targets"`
	Active         *bool             `json:"active"`
	Headers        map[string]string `json:"headers"`
	RateLimit      *RateLimit        `json:"rate_limit"`
}

// LoadBalancer represents a load balancer configuration
type LoadBalancer struct {
	ID        string                 `json:"id" db:"id"`
	Name      string                 `json:"name" db:"name"`
	Type      string                 `json:"type" db:"type"` // round-robin, least-connections, etc.
	Config    map[string]interface{} `json:"config" db:"config"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// LoadBalancerCreationRequest represents a request to create a new load balancer
type LoadBalancerCreationRequest struct {
	Name   string                 `json:"name" binding:"required"`
	Type   string                 `json:"type" binding:"required"`
	Config map[string]interface{} `json:"config"`
}
