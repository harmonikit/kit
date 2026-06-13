// Package etcd provides Etcd-based service discovery for harmoni/sd.
//
// This is a stub implementation. A production version depends on
// go.etcd.io/etcd/client/v3 for real registration and discovery.
package etcd

import (
	"context"
	"sync"

	"github.com/harmonikit/harmoni/sd"
)

// Compile-time interface checks.
var (
	_ sd.Registrar  = (*Registrar)(nil)
	_ sd.Discoverer = (*Discoverer)(nil)
	_ sd.Instancer  = (*Instancer)(nil)
)

// Registrar registers services with Etcd.
type Registrar struct {
	mu       sync.Mutex
	services []string
}

// NewRegistrar returns an Etcd Registrar.
func NewRegistrar() *Registrar {
	return &Registrar{}
}

func (r *Registrar) Register(_ context.Context) error {
	return nil
}

func (r *Registrar) Deregister(_ context.Context) error {
	return nil
}

// Discoverer discovers services from Etcd.
type Discoverer struct {
	instances []string
}

// NewDiscoverer returns a Discoverer with the given instances.
func NewDiscoverer(instances ...string) *Discoverer {
	return &Discoverer{instances: instances}
}

func (d *Discoverer) Discover(_ context.Context) ([]string, error) {
	result := make([]string, len(d.instances))
	copy(result, d.instances)
	return result, nil
}

// Instancer provides Etcd-based watch on instances.
type Instancer struct {
	mu        sync.Mutex
	instances []string
	subs      []chan struct{}
}

// NewInstancer returns an Instancer.
func NewInstancer(instances ...string) *Instancer {
	return &Instancer{instances: instances}
}

func (i *Instancer) Discover(ctx context.Context) ([]string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	result := make([]string, len(i.instances))
	copy(result, i.instances)
	return result, nil
}

func (i *Instancer) Subscribe() <-chan struct{} {
	i.mu.Lock()
	defer i.mu.Unlock()
	ch := make(chan struct{}, 1)
	i.subs = append(i.subs, ch)
	return ch
}

func (i *Instancer) Instances() ([]string, error) {
	return i.Discover(context.Background())
}

func (i *Instancer) Set(instances []string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.instances = instances
	for _, ch := range i.subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
