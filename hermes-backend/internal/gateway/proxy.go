package gateway

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
)

// Proxy handles reverse proxying requests to backend services
type Proxy struct {
	routes      map[string]*models.Route
	routesMutex sync.RWMutex
	log         logger.Logger
}

// NewProxy creates a new reverse proxy
func NewProxy(log logger.Logger) *Proxy {
	return &Proxy{
		routes: make(map[string]*models.Route),
		log:    log,
	}
}

// UpdateRoutes updates the proxy routes
func (p *Proxy) UpdateRoutes(routes []*models.Route) {
	p.routesMutex.Lock()
	defer p.routesMutex.Unlock()

	// Reset routes
	p.routes = make(map[string]*models.Route)

	// Add new routes
	for _, route := range routes {
		p.routes[route.ID] = route
	}

	p.log.Info().Msgf("Updated proxy routes, total routes: %d", len(routes))
}

// Handler returns an HTTP handler for the proxy
func (p *Proxy) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.routesMutex.RLock()
		defer p.routesMutex.RUnlock()

		// Find matching route
		var matchedRoute *models.Route
		for _, route := range p.routes {
			if strings.HasPrefix(r.URL.Path, route.Path) {
				matchedRoute = route
				break
			}
		}

		if matchedRoute == nil {
			http.Error(w, "Route not found", http.StatusNotFound)
			return
		}

		// Choose target based on load balancing strategy
		targetURL, err := p.chooseTarget(matchedRoute)
		if err != nil {
			p.log.Error().Err(err).Str("path", r.URL.Path).Msg("Failed to choose target")
			http.Error(w, "Failed to route request", http.StatusInternalServerError)
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Update request URL
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Header.Set("X-Forwarded-Proto", r.URL.Scheme)
		r.Header.Set("X-Forwarded-For", r.RemoteAddr)

		// Serve request
		proxy.ServeHTTP(w, r)
	})
}

// chooseTarget selects a target URL based on the route's load balancing strategy
func (p *Proxy) chooseTarget(route *models.Route) (*url.URL, error) {
	if len(route.Targets) == 0 {
		return nil, fmt.Errorf("no targets available for route: %s", route.Path)
	}

	// Simple round-robin for now
	// This would be replaced with more sophisticated load balancing in a real implementation
	target := route.Targets[0]

	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}

	return targetURL, nil
}
