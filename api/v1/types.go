package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// TenancyMode describes whether a sandbox can share a host or must have it exclusively.
type TenancyMode string

const (
	TenancyModeShared    TenancyMode = "shared"
	TenancyModeDedicated TenancyMode = "dedicated"
)

type Phase string

const (
	PhasePending       Phase = "Pending"
	PhaseProvisioning  Phase = "Provisioning"
	PhaseBootstrapping Phase = "Bootstrapping"
	PhaseRegistered    Phase = "Registered"
	PhaseReady         Phase = "Ready"
	PhaseScheduled     Phase = "Scheduled"
	PhaseCreating      Phase = "Creating"
	PhaseDraining      Phase = "Draining"
	PhaseReleasing     Phase = "Releasing"
	PhaseReleased      Phase = "Released"
	PhaseActive        Phase = "Active"
	PhaseDeleted       Phase = "Deleted"
	PhaseFailed        Phase = "Failed"
)

// ProviderRef identifies the external capacity source a pool or machine should use.
//
// Examples include Metal3, AWS EC2 bare metal, AWS EC2 Mac, or Scaleway.
type ProviderRef struct {
	Kind   string            `json:"kind,omitempty"`
	Region string            `json:"region,omitempty"`
	Zone   string            `json:"zone,omitempty"`
	Config map[string]string `json:"config,omitempty"`
}

// HostPolicy captures scheduling rules that constrain how sandboxes may share a host.
type HostPolicy struct {
	MaxActiveVMsPerHost int32 `json:"maxActiveVMsPerHost,omitempty"`
}

// WarmPoolSpec defines how much idle capacity Kubebox should try to keep ready.
type WarmPoolSpec struct {
	MinIdleHosts int32 `json:"minIdleHosts,omitempty"`
	MinIdleSlots int32 `json:"minIdleSlots,omitempty"`
	MaxHosts     int32 `json:"maxHosts,omitempty"`
}

// ScaleDownPolicy defines when a host may be released back to the underlying provider.
type ScaleDownPolicy struct {
	ReleasePolicy string `json:"releasePolicy,omitempty"`
}

// HostPoolSpec is the desired state of a host pool.
//
// In Kubernetes, "Spec" means "what the user wants."
type HostPoolSpec struct {
	Provider     ProviderRef   `json:"provider,omitempty"`
	Backend      string        `json:"backend,omitempty"`
	OS           string        `json:"os,omitempty"`
	Arch         string        `json:"arch,omitempty"`
	TenancyModes []TenancyMode `json:"tenancyModes,omitempty"`
	HostPolicy   HostPolicy    `json:"hostPolicy,omitempty"`
	WarmPool     WarmPoolSpec  `json:"warmPool,omitempty"`
	ScaleDown    ScaleDownPolicy `json:"scaleDown,omitempty"`
	Images       []string      `json:"images,omitempty"`
}

// HostPoolStatus is the observed state reported by controllers.
//
// In Kubernetes, "Status" means "what the system believes is currently true."
type HostPoolStatus struct {
	Phase              Phase              `json:"phase,omitempty"`
	DesiredHosts       int32              `json:"desiredHosts,omitempty"`
	ReadyHosts         int32              `json:"readyHosts,omitempty"`
	AvailableSlots     int32              `json:"availableSlots,omitempty"`
	LastScaleTime      *metav1.Time       `json:"lastScaleTime,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// ProviderMachineSpec describes one requested machine from an external provider.
//
// This is the bridge between Kubebox's desired capacity and a concrete machine
// requested from AWS, Scaleway, or another provider.
type ProviderMachineSpec struct {
	PoolRef   string      `json:"poolRef,omitempty"`
	Provider  ProviderRef `json:"provider,omitempty"`
	Image     string      `json:"image,omitempty"`
	Bootstrap string      `json:"bootstrap,omitempty"`
}

// ProviderMachineStatus reports what happened to that machine request after
// reconciliation against the external provider.
type ProviderMachineStatus struct {
	Phase          Phase              `json:"phase,omitempty"`
	ProviderID     string             `json:"providerID,omitempty"`
	ObservedHost   string             `json:"observedHost,omitempty"`
	LastUpdateTime *metav1.Time       `json:"lastUpdateTime,omitempty"`
	Conditions     []metav1.Condition `json:"conditions,omitempty"`
}

// HostSpec identifies a registered machine that the scheduler can reason about.
//
// A Host usually appears after a ProviderMachine has booted and the host agent
// has registered itself with the control plane.
type HostSpec struct {
	PoolRef            string `json:"poolRef,omitempty"`
	ProviderMachineRef string `json:"providerMachineRef,omitempty"`
	Backend            string `json:"backend,omitempty"`
	OS                 string `json:"os,omitempty"`
	Arch               string `json:"arch,omitempty"`
}

// HostStatus describes whether the machine is healthy and how much allocatable
// capacity it currently exposes to the scheduler.
type HostStatus struct {
	Phase            Phase              `json:"phase,omitempty"`
	Healthy          bool               `json:"healthy,omitempty"`
	AllocatableSlots int32              `json:"allocatableSlots,omitempty"`
	CachedImages     []string           `json:"cachedImages,omitempty"`
	LastHeartbeat    *metav1.Time       `json:"lastHeartbeat,omitempty"`
	Conditions       []metav1.Condition `json:"conditions,omitempty"`
}

// HostLeaseSpec is an exclusive claim on a host for a sandbox.
//
// Dedicated placement and one-active-VM-per-host policies are enforced through
// HostLease creation.
type HostLeaseSpec struct {
	HostRef    string `json:"hostRef,omitempty"`
	SandboxRef string `json:"sandboxRef,omitempty"`
}

// HostLeaseStatus reports whether the lease is active and when it expires.
type HostLeaseStatus struct {
	Phase      Phase        `json:"phase,omitempty"`
	ExpiryTime *metav1.Time `json:"expiryTime,omitempty"`
}

// ResourceRequirements captures the VM shape requested for a sandbox.
type ResourceRequirements struct {
	VCPU     int32 `json:"vcpu,omitempty"`
	MemoryMB int32 `json:"memoryMb,omitempty"`
	DiskGB   int32 `json:"diskGb,omitempty"`
}

// AccessPolicy controls which access channels Kubebox should publish for a sandbox.
type AccessPolicy struct {
	AllowSSH bool `json:"allowSSH,omitempty"`
	AllowGUI bool `json:"allowGUI,omitempty"`
}

// SandboxClassSpec defines a reusable sandbox template that products can refer to.
//
// It lets callers request "a kind of sandbox" without repeating all runtime and
// resource settings every time.
type SandboxClassSpec struct {
	Backend    string               `json:"backend,omitempty"`
	Tenancy    TenancyMode          `json:"tenancy,omitempty"`
	HostPolicy HostPolicy           `json:"hostPolicy,omitempty"`
	Resources  ResourceRequirements `json:"resources,omitempty"`
	Access     AccessPolicy         `json:"access,omitempty"`
	TTLSeconds int32                `json:"ttlSeconds,omitempty"`
}

// SandboxClassStatus reports the controller's view of the current template revision.
type SandboxClassStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// SandboxSpec is the desired state of one ephemeral VM sandbox.
//
// This is the main resource a product or higher-level API would create to ask
// Kubebox for an environment.
type SandboxSpec struct {
	Tenant     string               `json:"tenant,omitempty"`
	ClassRef   string               `json:"classRef,omitempty"`
	Image      string               `json:"image,omitempty"`
	Backend    string               `json:"backend,omitempty"`
	Tenancy    TenancyMode          `json:"tenancy,omitempty"`
	HostPolicy HostPolicy           `json:"hostPolicy,omitempty"`
	Resources  ResourceRequirements `json:"resources,omitempty"`
	Access     AccessPolicy         `json:"access,omitempty"`
	TTLSeconds int32                `json:"ttlSeconds,omitempty"`
}

// SandboxStatus reports where the sandbox landed and whether it is usable.
type SandboxStatus struct {
	Phase        Phase              `json:"phase,omitempty"`
	HostRef      string             `json:"hostRef,omitempty"`
	LeaseRef     string             `json:"leaseRef,omitempty"`
	AccessURL    string             `json:"accessURL,omitempty"`
	LastReadyAt  *metav1.Time       `json:"lastReadyAt,omitempty"`
	Conditions   []metav1.Condition `json:"conditions,omitempty"`
}

// HostPool is a Kubernetes custom resource.
//
// TypeMeta carries the API identity, like apiVersion and kind.
// ObjectMeta carries standard Kubernetes metadata, like name, namespace, labels,
// annotations, and resourceVersion.
//
// The usual pattern for Kubernetes resources is:
//   - Spec: desired state
//   - Status: observed state
//
// A HostPool declares a class of interchangeable hosts plus the warm-capacity
// policy Kubebox should maintain for them.
type HostPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HostPoolSpec   `json:"spec,omitempty"`
	Status HostPoolStatus `json:"status,omitempty"`
}

// HostPoolList is the list variant Kubernetes uses for list/watch operations.
type HostPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HostPool `json:"items"`
}

// ProviderMachine represents one concrete bare-metal machine request to an
// external provider.
type ProviderMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderMachineSpec   `json:"spec,omitempty"`
	Status ProviderMachineStatus `json:"status,omitempty"`
}

// ProviderMachineList is the list variant for ProviderMachine resources.
type ProviderMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderMachine `json:"items"`
}

// Host represents one registered machine that is available, or becoming
// available, for sandbox placement.
type Host struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HostSpec   `json:"spec,omitempty"`
	Status HostStatus `json:"status,omitempty"`
}

// HostList is the list variant for Host resources.
type HostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Host `json:"items"`
}

// HostLease represents exclusive ownership of a host by one sandbox.
type HostLease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HostLeaseSpec   `json:"spec,omitempty"`
	Status HostLeaseStatus `json:"status,omitempty"`
}

// HostLeaseList is the list variant for HostLease resources.
type HostLeaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HostLease `json:"items"`
}

// SandboxClass is a reusable template describing a family of sandboxes.
type SandboxClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SandboxClassSpec   `json:"spec,omitempty"`
	Status SandboxClassStatus `json:"status,omitempty"`
}

// SandboxClassList is the list variant for SandboxClass resources.
type SandboxClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SandboxClass `json:"items"`
}

// Sandbox is one requested ephemeral VM environment.
type Sandbox struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SandboxSpec   `json:"spec,omitempty"`
	Status SandboxStatus `json:"status,omitempty"`
}

// SandboxList is the list variant for Sandbox resources.
type SandboxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Sandbox `json:"items"`
}
