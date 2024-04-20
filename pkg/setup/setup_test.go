package setup_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/spf13/afero"
	"github.com/tomasaschan/envtest-at-runtime/pkg/setup"
)

func TestHappyPath(t *testing.T) {
	version := "1.27.x!"

	p, err := setup.EnsureBinariesExist(version,
		setup.WithLog(logr.FromSlogHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
		setup.WithFS(afero.NewMemMapFs()),
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		t.FailNow()
	}

	if p == "" {
		t.Error("Expected non-empty path")
	}

	t.Logf("passed! binaries at %s", p)
}

func TestInvalidVersion(t *testing.T) {
	version := "foo"

	p, err := setup.EnsureBinariesExist(version,
		setup.WithLog(logr.FromSlogHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
		setup.WithFS(afero.NewMemMapFs()),
	)

	if err == nil {
		t.Error("Expected error, got nil")
		t.Logf("path: %s", p)
		t.FailNow()
	}

	t.Logf("error: %v", err)
}
