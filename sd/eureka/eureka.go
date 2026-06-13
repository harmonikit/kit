// Package eureka provides Netflix Eureka-based service discovery for harmoni/sd.
//
// This is a stub implementation for service registration and discovery
// via a Eureka server.
package eureka

import (
	"context"
	"sync"

	"github.com/harmonikit/harmoni/sd"
)

// Compile-time interface checks.
var (
	_ sd.Registrar  = (*Registrar)(nil)
	_ sd.Discoverer = (*Discoverer)(nil)
)

// Registrar registers services with Eureka.
type Registrar struct {
	mu       sync.Mutex
	services []string
}

// NewRegistrar returns a Eureka Registrar.
func NewRegistrar() *Registrar {
	return &Registrar{mu: sync.Mutex{}, services: nil}
}

// Register adds the service instance.
func (r *Registrar) Register(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return nil
}

// Deregister removes the service instance.
func (r *Registrar) Deregister(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return nil
}

// Discoverer discovers services from Eureka.
type Discoverer struct {
	instances []string
}

// NewDiscoverer returns a Discoverer with the given instances.
func NewDiscoverer(instances ...string) *Discoverer {
	return &Discoverer{instances: instances}
}

// Discover returns the current set of instances.
func (d *Discoverer) Discover(_ context.Context) ([]string, error) {
	result := make([]string, len(d.instances))
	copy(result, d.instances)
	return result, nil
}
