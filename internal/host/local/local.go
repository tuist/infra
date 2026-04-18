package local

import runtimepkg "github.com/tuist/kubebox/internal/host/runtime"

type Runtime struct {
	hostID string
}

func New(hostID string) runtimepkg.Runtime {
	return &Runtime{hostID: hostID}
}

func (r *Runtime) Name() string {
	return "local"
}

func (r *Runtime) Capabilities() runtimepkg.Capabilities {
	return runtimepkg.Capabilities{
		HostID: r.hostID,
		Mode:   "local",
		Backends: []runtimepkg.BackendCapability{
			{
				Name:                "local",
				GuestOSes:           []string{"linux"},
				MaxActiveVMsPerHost: 1,
				SupportsShared:      false,
			},
		},
		OS:   "local",
		Arch: "amd64",
	}
}
