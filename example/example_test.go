package example_test

import (
	"context"
	"testing"

	"github.com/tomasaschan/envtest-at-runtime/pkg/setup"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestExample(t *testing.T) {
	path, err := setup.EnsureBinariesExist("1.27.x!")
	if err != nil {
		t.Fatalf("Failed to ensure envtest binaries are available: %v", err)
	}

	env := envtest.Environment{
		BinaryAssetsDirectory: path,
	}

	cfg, err := env.Start()
	if err != nil {
		t.Fatalf("Failed to start environment: %v", err)
	}

	c, err := client.New(cfg, client.Options{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.List(context.Background(), &v1.NamespaceList{}, &client.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to make a request to the test environment: %v", err)
	}
}
