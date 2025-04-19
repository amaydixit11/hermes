package gateway

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
)

// LoadBalancerType defines the type of load balancing strategy
type LoadBalancerType string

const (
	// RoundRobin distributes requests sequentially among targets
	RoundRobin LoadBalancerType = "round-robin"

	// LeastConnections routes to the target with the fewest active connections
	LeastConnections LoadBalancerType = "least-connections"

	// Random randomly selects a target
	Random LoadBalancerType = "random"

	// WeightedRoundRobin uses weights to determine distribution
	WeightedRoundRobin LoadBalancerType = "weighted-round-robin"
)

// LoadBalancer defines the interface for load balancing strategies
type LoadBalancer interface {
	// NextTarget returns the next target according to the load balancing strategy
	NextTarget(targets []string) (string, error)

	// Type returns the type of load balancer
	Type() LoadBalancerType
}

// ErrNoTargets is returned when no targets are available
var ErrNoTargets = errors.New("no targets available")

// RoundRobinBalancer implements a round-robin load balancing strategy
type RoundRobinBalancer struct {
	current uint64
}

// NewRoundRobinBalancer creates a new round-robin load balancer
func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		current: 0,
	}
}

// NextTarget returns the next target in round-robin fashion
func (lb *RoundRobinBalancer) NextTarget(targets []string) (string, error) {
	if len(targets) == 0 {
		return "", ErrNoTargets
	}

	next := atomic.AddUint64(&lb.current, 1) - 1
	return targets[next%uint64(len(targets))], nil
}

// Type returns the type of load balancer
func (lb *RoundRobinBalancer) Type() LoadBalancerType {
	return RoundRobin
}

// RandomBalancer implements a random load balancing strategy
type RandomBalancer struct {
	mu sync.Mutex
}

// NewRandomBalancer creates a new random load balancer
func NewRandomBalancer() *RandomBalancer {
	return &RandomBalancer{}
}

// NextTarget returns a randomly selected target
func (lb *RandomBalancer) NextTarget(targets []string) (string, error) {
	if len(targets) == 0 {
		return "", ErrNoTargets
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	idx := rand.Intn(len(targets))
	return targets[idx], nil
}

// Type returns the type of load balancer
func (lb *RandomBalancer) Type() LoadBalancerType {
	return Random
}

// LeastConnectionsBalancer implements a least-connections load balancing strategy
type LeastConnectionsBalancer struct {
	connections map[string]int
	mu          sync.Mutex
}

// NewLeastConnectionsBalancer creates a new least-connections load balancer
func NewLeastConnectionsBalancer() *LeastConnectionsBalancer {
	return &LeastConnectionsBalancer{
		connections: make(map[string]int),
	}
}

// NextTarget returns the target with the least active connections
func (lb *LeastConnectionsBalancer) NextTarget(targets []string) (string, error) {
	if len(targets) == 0 {
		return "", ErrNoTargets
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Initialize connections for new targets
	for _, target := range targets {
		if _, exists := lb.connections[target]; !exists {
			lb.connections[target] = 0
		}
	}

	// Find target with fewest connections
	var minTarget string
	minConnections := -1

	for _, target := range targets {
		connections := lb.connections[target]
		if minConnections == -1 || connections < minConnections {
			minConnections = connections
			minTarget = target
		}
	}

	// Increment connection count for selected target
	lb.connections[minTarget]++

	return minTarget, nil
}

// Type returns the type of load balancer
func (lb *LeastConnectionsBalancer) Type() LoadBalancerType {
	return LeastConnections
}

// ReleaseConnection decrements the connection count for a target
func (lb *LeastConnectionsBalancer) ReleaseConnection(target string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if count, exists := lb.connections[target]; exists && count > 0 {
		lb.connections[target]--
	}
}

// LoadBalancerFactory creates load balancers based on their type
type LoadBalancerFactory struct{}

// NewLoadBalancer creates a new load balancer of the specified type
func (f *LoadBalancerFactory) NewLoadBalancer(lbType LoadBalancerType) (LoadBalancer, error) {
	switch lbType {
	case RoundRobin:
		return NewRoundRobinBalancer(), nil
	case Random:
		return NewRandomBalancer(), nil
	case LeastConnections:
		return NewLeastConnectionsBalancer(), nil
	default:
		return nil, errors.New("unsupported load balancer type")
	}
}
