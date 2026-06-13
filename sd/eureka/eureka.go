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
	return &Registrar{}
}

func (r *Registrar) Register(ctx context.Context) error   { return nil }
func (r *Registrar) Deregister(ctx context.Context) error { return nil }

// Discoverer discovers services from Eureka.
type Discoverer struct {
	instances []string
}

// NewDiscoverer returns a Discoverer with the given instances.
func NewDiscoverer(instances ...string) *Discoverer {
	return &Discoverer{instances: instances}
}

func (d *Discoverer) Discover(ctx context.Context) ([]string, error) {
	result := make([]string, len(d.instances))
	copy(result, d.instances)
	return result, nil
}
