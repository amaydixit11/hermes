package repository

import (
	"context"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
)

// RouteRepository defines the interface for route data access
type RouteRepository interface {
	// Create creates a new route
	Create(ctx context.Context, route *models.Route) error

	// Get retrieves a route by ID
	Get(ctx context.Context, id string) (*models.Route, error)

	// List retrieves all routes with optional filtering
	List(ctx context.Context, filters map[string]interface{}) ([]*models.Route, error)

	// Update updates an existing route
	Update(ctx context.Context, id string, updates *models.RouteUpdateRequest) error

	// Delete removes a route
	Delete(ctx context.Context, id string) error
}

// LoadBalancerRepository defines the interface for load balancer data access
type LoadBalancerRepository interface {
	// Create creates a new load balancer
	Create(ctx context.Context, lb *models.LoadBalancer) error

	// Get retrieves a load balancer by ID
	Get(ctx context.Context, id string) (*models.LoadBalancer, error)

	// List retrieves all load balancers with optional filtering
	List(ctx context.Context, filters map[string]interface{}) ([]*models.LoadBalancer, error)

	// Update updates an existing load balancer
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete removes a load balancer
	Delete(ctx context.Context, id string) error
}
