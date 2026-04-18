package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeGroupVersion is the Kubernetes API identity for these custom resources.
	//
	// In YAML this becomes:
	//   apiVersion: kubebox.dev/v1
	//
	// "kubebox.dev" is the API group and "v1" is the API version.
	SchemeGroupVersion = schema.GroupVersion{Group: "kubebox.dev", Version: "v1"}

	// SchemeBuilder is a small controller-runtime helper used to register our Go
	// structs as Kubernetes API types.
	//
	// The leading '&' means "take the address of this struct", so this is a
	// pointer to one shared Builder instance rather than a copied value.
	SchemeBuilder      = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	// AddToScheme is the function the rest of the application calls when it wants
	// the Kubernetes runtime to recognize Kubebox resources such as HostPool and
	// Sandbox.
	AddToScheme        = SchemeBuilder.AddToScheme
)

// Resource returns a fully-qualified resource name in the current API group.
//
// Example:
//   Resource("hostpools") -> hostpools.kubebox.dev
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func init() {
	// Register tells the runtime which Go structs belong to this API version.
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
		&HostLease{},
		&HostLeaseList{},
		&SandboxClass{},
		&SandboxClassList{},
		&Sandbox{},
		&SandboxList{},
	)
}
