package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tuist/infra/internal/host/local"
	runtimepkg "github.com/tuist/infra/internal/host/runtime"
)

type Config struct {
	ListenAddress string
	Mode          string
	VMRuntimes    []string
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
				HostID:     cfg.HostID,
				Mode:       cfg.Mode,
				VMRuntimes: buildVMRuntimes(cfg.Mode, cfg.VMRuntimes),
				OS:         "linux",
				Arch:       "amd64",
			},
		}, nil
	case "macos":
		return &noopRuntime{
			name: "macos",
			capabilities: runtimepkg.Capabilities{
				HostID:     cfg.HostID,
				Mode:       cfg.Mode,
				VMRuntimes: buildVMRuntimes(cfg.Mode, cfg.VMRuntimes),
				OS:         "macos",
				Arch:       "arm64",
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported mode %q", cfg.Mode)
	}
}

func buildVMRuntimes(mode string, vmRuntimes []string) []runtimepkg.VMRuntimeCapability {
	normalized := normalizeVMRuntimes(mode, vmRuntimes)
	capabilities := make([]runtimepkg.VMRuntimeCapability, 0, len(normalized))

	for _, vmRuntime := range normalized {
		capabilities = append(capabilities, capabilityForVMRuntime(mode, vmRuntime))
	}

	return capabilities
}

func normalizeVMRuntimes(mode string, vmRuntimes []string) []string {
	filtered := make([]string, 0, len(vmRuntimes))
	for _, vmRuntime := range vmRuntimes {
		vmRuntime = strings.TrimSpace(vmRuntime)
		if vmRuntime == "" {
			continue
		}
		filtered = append(filtered, vmRuntime)
	}

	if len(filtered) > 0 {
		return filtered
	}

	switch mode {
	case "linux":
		return []string{"cloud-hypervisor"}
	case "macos":
		return []string{"tart"}
	default:
		return []string{"local"}
	}
}

func capabilityForVMRuntime(mode, vmRuntime string) runtimepkg.VMRuntimeCapability {
	switch mode {
	case "linux":
		return runtimepkg.VMRuntimeCapability{
			Name:                vmRuntime,
			GuestOSes:           []string{"linux"},
			MaxActiveVMsPerHost: 10,
			SupportsShared:      true,
		}
	case "macos":
		return runtimepkg.VMRuntimeCapability{
			Name:                vmRuntime,
			GuestOSes:           []string{"linux", "macos"},
			MaxActiveVMsPerHost: 1,
			SupportsShared:      false,
		}
	default:
		return runtimepkg.VMRuntimeCapability{
			Name:                vmRuntime,
			GuestOSes:           []string{"linux"},
			MaxActiveVMsPerHost: 1,
			SupportsShared:      false,
		}
	}
}
