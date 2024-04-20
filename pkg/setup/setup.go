package setup

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/go-logr/logr"
	"github.com/spf13/afero"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/env"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/remote"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/store"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/versions"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/workflows"
)

type SetupOption func(*env.Env)

func WithLog(log logr.Logger) SetupOption {
	return func(env *env.Env) {
		env.Log = log
	}
}

func WithClient(client remote.Client) SetupOption {
	return func(env *env.Env) {
		env.Client = &client
	}
}

func WithPlatform(platform versions.PlatformItem) SetupOption {
	return func(env *env.Env) {
		env.Platform = platform
	}
}

func WithFS(fs afero.Fs) SetupOption {
	return func(env *env.Env) {
		env.FS = afero.Afero{Fs: fs}
	}
}

func WithStoreDir(dir string) SetupOption {
	return func(env *env.Env) {
		env.Store = store.NewAt(dir)
	}
}

func getConfig(version versions.Spec, out io.Writer, options ...SetupOption) (*env.Env, error) {
	config := &env.Env{
		Version: version,
		Platform: versions.PlatformItem{
			Platform: versions.Platform{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
		},
		Out: out,
		FS:  afero.Afero{Fs: afero.NewOsFs()},
	}

	for _, opt := range options {
		opt(config)
	}

	if config.Store == nil {
		binPath, err := store.DefaultStoreDir()
		if err != nil {
			return nil, fmt.Errorf("get default store dir: %w", err)
		}
		config.Store = store.NewAt(binPath)
	}

	if config.Log.IsZero() {
		config.Log = logr.FromSlogHandler(slog.NewTextHandler(os.Stdout, nil))
	}

	if config.Client == nil {
		config.Client = &remote.Client{
			Bucket: "kubebuilder-tools",
			Server: "storage.googleapis.com",
			Log:    config.Log.WithName("storage-client"),
		}
	}

	return config, nil
}

func EnsureBinariesExist(version string, options ...SetupOption) (string, error) {
	v, err := versions.FromExpr(version)
	if err != nil {
		return "", fmt.Errorf("invalid version %q: %w", version, err)
	}

	out := strings.Builder{}

	config, err := getConfig(v, &out, options...)

	var p any
	defer func() {
		p = recover()
	}()

	workflows.Use{PrintFormat: env.PrintPath}.Do(config)

	if p != nil {
		return "", fmt.Errorf("execute setup-envtest use: %v", err)
	}

	return out.String(), nil
}
