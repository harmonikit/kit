// Package dns provides DNS-based service discovery for harmoni/sd.
//
// It uses stdlib net.LookupSRV for discovering instances via DNS SRV records.
package dns

import (
	"context"
	"fmt"
	"net"

	"github.com/harmonikit/harmoni/sd"
)

// Compile-time interface checks.
var (
	_ sd.Discoverer = (*Discoverer)(nil)
)

// Discoverer discovers services via DNS SRV records.
type Discoverer struct {
	service string
	proto   string
	name    string
}

// NewDiscoverer creates a DNS discoverer for the given SRV record.
// For example: NewDiscoverer("service", "tcp", "example.com") will query
// _service._tcp.example.com.
func NewDiscoverer(service, proto, name string) *Discoverer {
	return &Discoverer{service: service, proto: proto, name: name}
}

// Discover performs a DNS SRV lookup and returns instance addresses.
func (d *Discoverer) Discover(ctx context.Context) ([]string, error) {
	_, addrs, err := net.DefaultResolver.LookupSRV(ctx, d.service, d.proto, d.name)
	if err != nil {
		return nil, fmt.Errorf("dns srv lookup: %w", err)
	}
	instances := make([]string, len(addrs))
	for i, addr := range addrs {
		instances[i] = fmt.Sprintf("%s:%d", addr.Target, addr.Port)
	}
	return instances, nil
}
