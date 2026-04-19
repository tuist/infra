package controller

import infrav1 "github.com/tuist/infra/api/v1"

// poolSupportsBackend reports whether a HostPool is meant to supply hosts for
// a given runtime backend.
//
// Example:
//
//	pool.spec.backends = ["cloud-hypervisor", "tart"]
//	backend = "tart"
//	-> true
func poolSupportsBackend(pool infrav1.HostPool, backend string) bool {
	for _, candidate := range pool.Spec.Backends {
		if candidate == backend {
			return true
		}
	}

	return false
}

// hostSupportsSandbox reports whether a specific host advertises a backend
// capability that can run the requested sandbox.
//
// The scheduler first matches the backend name, then optionally matches the
// guest operating system if the sandbox asked for one.
//
// Example:
//
//	host.spec.backends = [{name: "tart", guestOSes: ["linux", "macos"]}]
//	sandbox.spec.backend = "tart"
//	sandbox.spec.guestOS = "linux"
//	-> true
//
// Example:
//
//	host.spec.backends = [{name: "tart", guestOSes: ["linux", "macos"]}]
//	sandbox.spec.backend = "cloud-hypervisor"
//	-> false
func hostSupportsSandbox(host infrav1.Host, sandbox infrav1.Sandbox) bool {
	for _, capability := range host.Spec.Backends {
		if capability.Name != sandbox.Spec.Backend {
			continue
		}

		if sandbox.Spec.GuestOS == "" || contains(capability.GuestOSes, sandbox.Spec.GuestOS) {
			return true
		}
	}

	return false
}

// contains is a small helper for matching one requested value against a list of
// values advertised by a host capability.
func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
