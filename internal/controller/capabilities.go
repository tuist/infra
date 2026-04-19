package controller

import infrav1 "github.com/tuist/infra/api/v1"

// poolSupportsVMRuntime reports whether a HostPool is meant to supply hosts for
// a given VM runtime.
//
// Example:
//
//	pool.spec.vmRuntimes = ["cloud-hypervisor", "tart"]
//	vmRuntime = "tart"
//	-> true
func poolSupportsVMRuntime(pool infrav1.HostPool, vmRuntime string) bool {
	for _, candidate := range pool.Spec.VMRuntimes {
		if candidate == vmRuntime {
			return true
		}
	}

	return false
}

// hostSupportsSandbox reports whether a specific host advertises a VM runtime
// capability that can run the requested sandbox.
//
// The scheduler first matches the VM runtime name, then optionally matches the
// guest operating system if the sandbox asked for one.
//
// Example:
//
//	host.spec.vmRuntimes = [{name: "tart", guestOSes: ["linux", "macos"]}]
//	sandbox.spec.vmRuntime = "tart"
//	sandbox.spec.guestOS = "linux"
//	-> true
//
// Example:
//
//	host.spec.vmRuntimes = [{name: "tart", guestOSes: ["linux", "macos"]}]
//	sandbox.spec.vmRuntime = "cloud-hypervisor"
//	-> false
func hostSupportsSandbox(host infrav1.Host, sandbox infrav1.Sandbox) bool {
	for _, capability := range host.Spec.VMRuntimes {
		if capability.Name != sandbox.Spec.VMRuntime {
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
