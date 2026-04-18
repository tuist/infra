package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tuist/kubebox/internal/host/local"
	runtimepkg "github.com/tuist/kubebox/internal/host/runtime"
)

type Config struct {
	ListenAddress string
	Mode          string
	Backend       string
	HostID        string
}

type Agent struct {
	cfg        Config
	runtime    runtimepkg.Runtime
	httpServer *http.Server
}

type noopRuntime struct {
	name         string
	capabilities runtimepkg.Capabilities
}

func (r *noopRuntime) Name() string {
	return r.name
}

func (r *noopRuntime) Capabilities() runtimepkg.Capabilities {
	return r.capabilities
}

func New(cfg Config) (*Agent, error) {
	rt, err := newRuntime(cfg)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/v1/state", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"runtime":      rt.Name(),
			"capabilities": rt.Capabilities(),
		})
	})

	return &Agent{
		cfg:     cfg,
		runtime: rt,
		httpServer: &http.Server{
			Addr:              cfg.ListenAddress,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}, nil
}

func (a *Agent) ListenAndServe() error {
	return a.httpServer.ListenAndServe()
}

func newRuntime(cfg Config) (runtimepkg.Runtime, error) {
	if cfg.HostID == "" {
		return nil, fmt.Errorf("host ID is required")
	}

	switch cfg.Mode {
	case "local":
		return local.New(cfg.HostID), nil
	case "linux":
		return &noopRuntime{
			name: "linux",
			capabilities: runtimepkg.Capabilities{
				HostID:              cfg.HostID,
				Mode:                cfg.Mode,
				Backend:             cfg.Backend,
				OS:                  "linux",
				Arch:                "amd64",
				MaxActiveVMsPerHost: 10,
				SupportsShared:      true,
			},
		}, nil
	case "macos":
		return &noopRuntime{
			name: "macos",
			capabilities: runtimepkg.Capabilities{
				HostID:              cfg.HostID,
				Mode:                cfg.Mode,
				Backend:             cfg.Backend,
				OS:                  "macos",
				Arch:                "arm64",
				MaxActiveVMsPerHost: 1,
				SupportsShared:      false,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported mode %q", cfg.Mode)
	}
}
