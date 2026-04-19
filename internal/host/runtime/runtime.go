package runtime

type VMRuntimeCapability struct {
	Name                string   `json:"name"`
	GuestOSes           []string `json:"guestOSes"`
	MaxActiveVMsPerHost int32    `json:"maxActiveVMsPerHost"`
	SupportsShared      bool     `json:"supportsShared"`
}

type Capabilities struct {
	HostID     string                `json:"hostID"`
	Mode       string                `json:"mode"`
	VMRuntimes []VMRuntimeCapability `json:"vmRuntimes"`
	OS         string                `json:"os"`
	Arch       string                `json:"arch"`
}

type Runtime interface {
	Name() string
	Capabilities() Capabilities
}
