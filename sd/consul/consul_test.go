package consul_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/sd"
	consulsd "github.com/harmonikit/kit/sd/consul"
)

func TestDiscoverer_Discover(t *testing.T) {
	d := consulsd.NewDiscoverer("host1:8080", "host2:8080")

	instances, err := d.Discover(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatalf("got %d, want 2", len(instances))
	}
}

func TestInstancer_Subscribe(t *testing.T) {
	i := consulsd.NewInstancer("host1:8080")

	ch := i.Subscribe()
	i.Set([]string{"host1:8080", "host2:8080"})

	select {
	case <-ch:
		// Notification received.
	default:
		t.Fatal("expected notification on instance change")
	}
}

func TestInstancer_Instances(t *testing.T) {
	i := consulsd.NewInstancer("a:1", "b:2")

	instances, err := i.Instances()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatalf("got %d, want 2", len(instances))
	}
}

func TestRegistrar(t *testing.T) {
	r := consulsd.NewRegistrar()
	if err := r.Register(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Deregister(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterfaces(t *testing.T) {
	var _ sd.Registrar = consulsd.NewRegistrar()
	var _ sd.Discoverer = consulsd.NewDiscoverer()
	var _ sd.Instancer = consulsd.NewInstancer()
}
