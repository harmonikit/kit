package dns_test

import (
	"testing"

	"github.com/harmonikit/harmoni/sd"
	dnssd "github.com/harmonikit/kit/sd/dns"
)

func TestDiscoverer_New(t *testing.T) {
	d := dnssd.NewDiscoverer("http", "tcp", "example.com")
	if d == nil {
		t.Fatal("expected non-nil discoverer")
	}
}

func TestDiscoverer_TypeCheck(t *testing.T) {
	var _ sd.Discoverer = dnssd.NewDiscoverer("s", "tcp", "n")
}
