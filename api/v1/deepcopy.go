package v1

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
)

func clone[T any](in *T) *T {
	if in == nil {
		return nil
	}

	out := new(T)
	payload, err := json.Marshal(in)
	if err != nil {
		*out = *in
		return out
	}

	if err := json.Unmarshal(payload, out); err != nil {
		*out = *in
	}

	return out
}

func (in *HostPool) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *HostPoolList) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *ProviderMachine) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *ProviderMachineList) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *Host) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *HostList) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *HostLease) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *HostLeaseList) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *SandboxClass) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *SandboxClassList) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *Sandbox) DeepCopyObject() runtime.Object {
	return clone(in)
}

func (in *SandboxList) DeepCopyObject() runtime.Object {
	return clone(in)
}
