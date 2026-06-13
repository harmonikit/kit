package etcd_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/sd"
	etcdsd "github.com/harmonikit/kit/sd/etcd"
)

func TestDiscoverer_Discover(t *testing.T) {
	d := etcdsd.NewDiscoverer("host1:9090", "host2:9090")

	instances, err := d.Discover(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatalf("got %d, want 2", len(instances))
	}
}

func TestInstancer_Subscribe(t *testing.T) {
	i := etcdsd.NewInstancer("host1:9090")

	ch := i.Subscribe()
	i.Set([]string{"host1:9090", "host2:9090"})

	select {
	case <-ch:
	default:
		t.Fatal("expected notification")
	}
}

func TestInstancer_Instances(t *testing.T) {
	i := etcdsd.NewInstancer("a:1")

	instances, err := i.Instances()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 1 {
		t.Fatalf("got %d, want 1", len(instances))
	}
}

func TestRegistrar(t *testing.T) {
	r := etcdsd.NewRegistrar()
	if err := r.Register(context.Background()); err != nil {
		t.Fatalf("unexpected register: %v", err)
	}
	if err := r.Deregister(context.Background()); err != nil {
		t.Fatalf("unexpected deregister: %v", err)
	}
}

func TestInterfaces(t *testing.T) {
	var _ sd.Registrar = etcdsd.NewRegistrar()
	var _ sd.Discoverer = etcdsd.NewDiscoverer()
	var _ sd.Instancer = etcdsd.NewInstancer()
}
