package runtime

type Capabilities struct {
	HostID              string `json:"hostID"`
	Mode                string `json:"mode"`
	Backend             string `json:"backend"`
	OS                  string `json:"os"`
	Arch                string `json:"arch"`
	MaxActiveVMsPerHost int32  `json:"maxActiveVMsPerHost"`
	SupportsShared      bool   `json:"supportsShared"`
}

type Runtime interface {
	Name() string
	Capabilities() Capabilities
}
