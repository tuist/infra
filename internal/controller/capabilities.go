package controller

import kubeboxv1 "github.com/tuist/kubebox/api/v1"

func poolSupportsBackend(pool kubeboxv1.HostPool, backend string) bool {
	for _, candidate := range pool.Spec.Backends {
		if candidate == backend {
			return true
		}
	}

	return false
}

func hostSupportsSandbox(host kubeboxv1.Host, sandbox kubeboxv1.Sandbox) bool {
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

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
