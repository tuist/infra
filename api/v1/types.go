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

type ProviderRef struct {
	Kind   string            `json:"kind,omitempty"`
	Region string            `json:"region,omitempty"`
	Zone   string            `json:"zone,omitempty"`
	Config map[string]string `json:"config,omitempty"`
}

type HostPolicy struct {
	MaxActiveVMsPerHost int32 `json:"maxActiveVMsPerHost,omitempty"`
}

type WarmPoolSpec struct {
	MinIdleHosts int32 `json:"minIdleHosts,omitempty"`
	MinIdleSlots int32 `json:"minIdleSlots,omitempty"`
	MaxHosts     int32 `json:"maxHosts,omitempty"`
}

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

type ProviderMachineSpec struct {
	PoolRef   string      `json:"poolRef,omitempty"`
	Provider  ProviderRef `json:"provider,omitempty"`
	Image     string      `json:"image,omitempty"`
	Bootstrap string      `json:"bootstrap,omitempty"`
}

type ProviderMachineStatus struct {
	Phase          Phase              `json:"phase,omitempty"`
	ProviderID     string             `json:"providerID,omitempty"`
	ObservedHost   string             `json:"observedHost,omitempty"`
	LastUpdateTime *metav1.Time       `json:"lastUpdateTime,omitempty"`
	Conditions     []metav1.Condition `json:"conditions,omitempty"`
}

type HostSpec struct {
	PoolRef            string `json:"poolRef,omitempty"`
	ProviderMachineRef string `json:"providerMachineRef,omitempty"`
	Backend            string `json:"backend,omitempty"`
	OS                 string `json:"os,omitempty"`
	Arch               string `json:"arch,omitempty"`
}

type HostStatus struct {
	Phase            Phase              `json:"phase,omitempty"`
	Healthy          bool               `json:"healthy,omitempty"`
	AllocatableSlots int32              `json:"allocatableSlots,omitempty"`
	CachedImages     []string           `json:"cachedImages,omitempty"`
	LastHeartbeat    *metav1.Time       `json:"lastHeartbeat,omitempty"`
	Conditions       []metav1.Condition `json:"conditions,omitempty"`
}

type HostLeaseSpec struct {
	HostRef    string `json:"hostRef,omitempty"`
	SandboxRef string `json:"sandboxRef,omitempty"`
}

type HostLeaseStatus struct {
	Phase      Phase        `json:"phase,omitempty"`
	ExpiryTime *metav1.Time `json:"expiryTime,omitempty"`
}

type ResourceRequirements struct {
	VCPU     int32 `json:"vcpu,omitempty"`
	MemoryMB int32 `json:"memoryMb,omitempty"`
	DiskGB   int32 `json:"diskGb,omitempty"`
}

type AccessPolicy struct {
	AllowSSH bool `json:"allowSSH,omitempty"`
	AllowGUI bool `json:"allowGUI,omitempty"`
}

type SandboxClassSpec struct {
	Backend    string               `json:"backend,omitempty"`
	Tenancy    TenancyMode          `json:"tenancy,omitempty"`
	HostPolicy HostPolicy           `json:"hostPolicy,omitempty"`
	Resources  ResourceRequirements `json:"resources,omitempty"`
	Access     AccessPolicy         `json:"access,omitempty"`
	TTLSeconds int32                `json:"ttlSeconds,omitempty"`
}

type SandboxClassStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

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
