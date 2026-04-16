package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	SchemeGroupVersion = schema.GroupVersion{Group: "kubebox.dev", Version: "v1"}
	SchemeBuilder      = &scheme.Builder{GroupVersion: SchemeGroupVersion}
	AddToScheme        = SchemeBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func init() {
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
