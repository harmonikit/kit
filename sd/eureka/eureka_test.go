package eureka_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/sd"
	eurekasd "github.com/harmonikit/kit/sd/eureka"
)

func TestDiscoverer_Discover(t *testing.T) {
	d := eurekasd.NewDiscoverer("host1:8761", "host2:8761")

	instances, err := d.Discover(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatalf("got %d, want 2", len(instances))
	}
}

func TestRegistrar(t *testing.T) {
	r := eurekasd.NewRegistrar()
	if err := r.Register(context.Background()); err != nil {
		t.Fatalf("unexpected register: %v", err)
	}
	if err := r.Deregister(context.Background()); err != nil {
		t.Fatalf("unexpected deregister: %v", err)
	}
}

func TestInterfaces(t *testing.T) {
	var _ sd.Registrar = eurekasd.NewRegistrar()
	var _ sd.Discoverer = eurekasd.NewDiscoverer()
}
