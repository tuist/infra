package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeGroupVersion is the Kubernetes API identity for these custom resources.
	//
	// In YAML this becomes:
	//   apiVersion: infra.tuist.dev/v1
	//
	// "infra.tuist.dev" is the API group and "v1" is the API version.
	SchemeGroupVersion = schema.GroupVersion{Group: "infra.tuist.dev", Version: "v1"}

	// SchemeBuilder is a small controller-runtime helper used to register our
	// Kubernetes API types under this group and version.
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	// AddToScheme is the function the rest of the application calls when it wants
	// the Kubernetes runtime to recognize Infra resources such as HostPool and
	// Sandbox.
	AddToScheme = SchemeBuilder.AddToScheme
)

// Resource returns a fully-qualified resource name in the current API group.
//
// Example:
//
//	Resource("hostpools") -> hostpools.infra.tuist.dev
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func init() {
	// Register tells the runtime which resource types belong to this API version.
	//
	// For every top-level resource Kubernetes also expects a corresponding List
	// type for list/watch operations.
	SchemeBuilder.Register(
		&HostPool{},
		&HostPoolList{},
		&ProviderMachine{},
		&ProviderMachineList{},
		&Host{},
		&HostList{},
		&HostImage{},
		&HostImageList{},
		&HostLease{},
		&HostLeaseList{},
		&SandboxClass{},
		&SandboxClassList{},
		&Sandbox{},
		&SandboxList{},
	)
}
