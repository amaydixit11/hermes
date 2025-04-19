// package postgres

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/jmoiron/sqlx"

// 	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
// 	"github.com/amaydixit11/hermes/hermes-backend/internal/repository"
// )

// // routeRepository is a PostgreSQL implementation of repository.RouteRepository
// type routeRepository struct {
// 	db *sqlx.DB
// }

// // NewRouteRepository creates a new PostgreSQL route repository
// func NewRouteRepository(db *sqlx.DB) repository.RouteRepository {
// 	return &routeRepository{db: db}
// }

// // Create creates a new route
// func (r *routeRepository) Create(ctx context.Context, route *models.Route) error {
// 	if route.ID == "" {
// 		route.ID = uuid.New().String()
// 	}

// 	now := time.Now()
// 	route.CreatedAt = now
// 	route.UpdatedAt = now

// 	targetsJSON, err := json.Marshal(route.Targets)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal targets: %w", err)
// 	}

// 	headersJSON, err := json.Marshal(route.Headers)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal headers: %w", err)
// 	}

// 	var rateLimitJSON []byte
// 	if route.RateLimit != nil {
// 		rateLimitJSON, err = json.Marshal(route.RateLimit)
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal rate limit: %w", err)
// 		}
// 	}

// 	query := `
// 		INSERT INTO routes (
// 			id, path, description, service_id, load_balancer_id,
// 			targets, active, headers, rate_limit, created_at, updated_at
// 		) VALUES (
// 			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
// 		)
// 	`

// 	_, err = r.db.ExecContext(
// 		ctx,
// 		query,
// 		route.ID,
// 		route.Path,
// 		route.Description,
// 		route.ServiceID,
// 		route.LoadBalancerID,
// 		targetsJSON,
// 		route.Active,
// 		headersJSON,
// 		rateLimitJSON,
// 		route.CreatedAt,
// 		route.UpdatedAt,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("failed to create route: %w", err)
// 	}

// 	return nil
// }

// // Get retrieves a route by ID
// func (r *routeRepository) Get(ctx context.Context, id string) (*models.Route, error) {
// 	query := `
// 		SELECT
// 			id, path, description, service_id, load_balancer_id,
// 			targets, active, headers, rate_limit, created_at, updated_at
// 		FROM routes
// 		WHERE id = $1
// 	`

// 	var (
// 		route         models.Route
// 		targetsJSON   []byte
// 		headersJSON   []byte
// 		rateLimitJSON []byte
// 	)

// 	err := r.db.QueryRowContext(
// 		ctx,
// 		query,
// 		id,
// 	).Scan(
// 		&route.ID,
// 		&route.Path,
// 		&route.Description,
// 		&route.ServiceID,
// 		&route.LoadBalancerID,
// 		&targetsJSON,
// 		&route.Active,
// 		&headersJSON,
// 		&rateLimitJSON,
// 		&route.CreatedAt,
// 		&route.UpdatedAt,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, repository.ErrNotFound
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get route: %w", err)
// 	}

// 	if err := json.Unmarshal(targetsJSON, &route.Targets); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal targets: %w", err)
// 	}

// 	if err := json.Unmarshal(headersJSON, &route.Headers); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal headers: %w", err)
// 	}

// 	if len(rateLimitJSON) > 0 {
// 		route.RateLimit = &models.RateLimit{}
// 		if err := json.Unmarshal(rateLimitJSON, route.RateLimit); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal rate limit: %w", err)
// 		}
// 	}

// 	return &route, nil
// }

// // List retrieves all routes with optional filtering
// func (r *routeRepository) List(ctx context.Context, filters map[string]interface{}) ([]*models.Route, error) {
// 	// Base query
// 	query := `
// 		SELECT
// 			id, path, description, service_id, load_balancer_id,
// 			targets, active, headers, rate_limit, created_at, updated_at
// 		FROM routes
// 	`

// 	// Apply filters if any
// 	var args []interface{}
// 	whereClause := ""

// 	if len(filters) > 0 {
// 		whereClause = " WHERE "
// 		i := 1

// 		for key, value := range filters {
// 			if i > 1 {
// 				whereClause += " AND "
// 			}

// 			whereClause += fmt.Sprintf("%s = $%d", key, i)
// 			args = append(args, value)
// 			i++
// 		}

// 		query += whereClause
// 	}

// 	rows, err := r.db.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list routes: %w", err)
// 	}
// 	defer rows.Close()

// 	var routes []*models.Route

// 	for rows.Next() {
// 		var (
// 			route         models.Route
// 			targetsJSON   []byte
// 			headersJSON   []byte
// 			rateLimitJSON []byte
// 		)

// 		err := rows.Scan(
// 			&route.ID,
// 			&route.Path,
// 			&route.Description,
// 			&route.ServiceID,
// 			&route.LoadBalancerID,
// 			&targetsJSON,
// 			&route.Active,
// 			&headersJSON,
// 			&rateLimitJSON,
// 			&route.CreatedAt,
// 			&route.UpdatedAt,
// 		)

// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan route: %w", err)
// 		}

// 		if err := json.Unmarshal(targetsJSON, &route.Targets); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal targets: %w", err)
// 		}

// 		if err := json.Unmarshal(headersJSON, &route.Headers); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal headers: %w", err)
// 		}

// 		if len(rateLimitJSON) > 0 {
// 			route.RateLimit = &models.RateLimit{}
// 			if err := json.Unmarshal(rateLimitJSON, route.RateLimit); err != nil {
// 				return nil, fmt.Errorf("failed to unmarshal rate limit: %w", err)
// 			}
// 		}

// 		routes = append(routes, &route)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating routes rows: %w", err)
// 	}

// 	return routes, nil
// }

// // Update updates an existing route
// func (r *routeRepository) Update(ctx context.Context, id string, updates *models.RouteUpdateRequest) error {
// 	// First, get the current route
// 	current, err := r.Get(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	// Apply updates
// 	if updates.Path != nil {
// 		current.Path = *updates.Path
// 	}

// 	if updates.Description != nil {
// 		current.Description = *updates.Description
// 	}

// 	if updates.ServiceID != nil {
// 		current.ServiceID = *updates.ServiceID
// 	}

// 	if updates.LoadBalancerID != nil {
// 		current.LoadBalancerID = *updates.LoadBalancerID
// 	}

// 	if updates.Targets != nil {
// 		current.Targets = updates.Targets
// 	}

// 	if updates.Active != nil {
// 		current.Active = *updates.Active
// 	}

// 	if updates.Headers != nil {
// 		current.Headers = updates.Headers
// 	}

// 	if updates.RateLimit != nil {
// 		current.RateLimit = updates.RateLimit
// 	}

// 	current.UpdatedAt = time.Now()

// 	// Convert to JSON
// 	targetsJSON, err := json.Marshal(current.Targets)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal targets: %w", err)
// 	}

// 	headersJSON, err := json.Marshal(current.Headers)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal headers: %w", err)
// 	}

// 	var rateLimitJSON []byte
// 	if current.RateLimit != nil {
// 		rateLimitJSON, err = json.Marshal(current.RateLimit)
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal rate limit: %w", err)
// 		}
// 	}

// 	query := `
// 		UPDATE routes
// 		SET
// 			path = $1,
// 			description = $2,
// 			service_id = $3,
// 			load_balancer_id = $4,
// 			targets = $5,
// 			active = $6,
// 			headers = $7,
// 			rate_limit = $8,
// 			updated_at = $9
// 		WHERE id = $10
// 	`

// 	_, err = r.db.ExecContext(
// 		ctx,
// 		query,
// 		current.Path,
// 		current.Description,
// 		current.ServiceID,
// 		current.LoadBalancerID,
// 		targetsJSON,
// 		current.Active,
// 		headersJSON,
// 		rateLimitJSON,
// 		current.UpdatedAt,
// 		id,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("failed to update route: %w", err)
// 	}

// 	return nil
// }

// // Delete removes a route
// func (r *routeRepository) Delete(ctx context.Context, id string) error {
// 	query := "DELETE FROM routes WHERE id = $1"

// 	result, err := r.db.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete route: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get rows affected: %w", err)
// 	}

// 	if rowsAffected == 0 {
// 		return repository.ErrNotFound
// 	}

// 	return nil
// }
