package dns_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/sd"
	dnssd "github.com/harmonikit/kit/sd/dns"
)

func TestDiscoverer_Interface(t *testing.T) {
	d := dnssd.NewDiscoverer("service", "tcp", "example.com")
	var _ sd.Discoverer = d
}

func TestDiscoverer_Discover_Invalid(t *testing.T) {
	d := dnssd.NewDiscoverer("nonexistent", "tcp", "invalid.local")

	ctx := context.Background()
	_, err := d.Discover(ctx)
	if err == nil {
		t.Fatal("expected error for invalid DNS lookup")
	}
}
