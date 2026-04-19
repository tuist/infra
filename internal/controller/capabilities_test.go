package controller

import (
	"testing"

	infrav1 "github.com/tuist/infra/api/v1"
)

func TestPoolSupportsVMRuntime(t *testing.T) {
	t.Parallel()

	pool := infrav1.HostPool{
		Spec: infrav1.HostPoolSpec{
			VMRuntimes: []string{"cloud-hypervisor", "tart"},
		},
	}

	if !poolSupportsVMRuntime(pool, "tart") {
		t.Fatal("expected pool to support tart")
	}

	if poolSupportsVMRuntime(pool, "firecracker") {
		t.Fatal("expected pool not to support firecracker")
	}
}

func TestHostSupportsSandbox(t *testing.T) {
	t.Parallel()

	host := infrav1.Host{
		Spec: infrav1.HostSpec{
			VMRuntimes: []infrav1.VMRuntimeCapability{
				{
					Name:      "cloud-hypervisor",
					GuestOSes: []string{"linux"},
				},
				{
					Name:      "tart",
					GuestOSes: []string{"linux", "macos"},
				},
			},
		},
	}

	testCases := []struct {
		name    string
		sandbox infrav1.Sandbox
		want    bool
	}{
		{
			name: "matching runtime and guest os",
			sandbox: infrav1.Sandbox{
				Spec: infrav1.SandboxSpec{
					VMRuntime: "tart",
					GuestOS:   "macos",
				},
			},
			want: true,
		},
		{
			name: "matching runtime without guest os",
			sandbox: infrav1.Sandbox{
				Spec: infrav1.SandboxSpec{
					VMRuntime: "tart",
				},
			},
			want: true,
		},
		{
			name: "matching runtime but unsupported guest os",
			sandbox: infrav1.Sandbox{
				Spec: infrav1.SandboxSpec{
					VMRuntime: "cloud-hypervisor",
					GuestOS:   "macos",
				},
			},
			want: false,
		},
		{
			name: "unsupported runtime",
			sandbox: infrav1.Sandbox{
				Spec: infrav1.SandboxSpec{
					VMRuntime: "qemu",
					GuestOS:   "linux",
				},
			},
			want: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := hostSupportsSandbox(host, testCase.sandbox)
			if got != testCase.want {
				t.Fatalf("hostSupportsSandbox() = %t, want %t", got, testCase.want)
			}
		})
	}
}

func TestHostHasReadyImageForSandbox(t *testing.T) {
	t.Parallel()

	hostImages := []infrav1.HostImage{
		{
			Spec: infrav1.HostImageSpec{
				HostRef:   "host-a",
				VMRuntime: "tart",
				Image:     "macos-14-base",
			},
			Status: infrav1.HostImageStatus{
				Phase: infrav1.PhaseReady,
			},
		},
		{
			Spec: infrav1.HostImageSpec{
				HostRef:   "host-a",
				VMRuntime: "cloud-hypervisor",
				Image:     "ubuntu-24.04",
			},
			Status: infrav1.HostImageStatus{
				Phase: infrav1.PhasePulling,
			},
		},
	}

	testCases := []struct {
		name     string
		hostName string
		sandbox  infrav1.Sandbox
		want     bool
	}{
		{
			name:     "no image requested",
			hostName: "host-a",
			sandbox: infrav1.Sandbox{
				Spec: infrav1.SandboxSpec{
					VMRuntime: "tart",
				},
			},
			want: true,
		},
		{
			name:     "ready image for host and runtime",
			hostName: "host-a",
			sandbox: infrav1.Sandbox{
				Spec: infrav1.SandboxSpec{
					VMRuntime: "tart",
					Image:     "macos-14-base",
				},
			},
			want: true,
		},
		{
			name:     "image still pulling",
			hostName: "host-a",
			sandbox: infrav1.Sandbox{
				Spec: infrav1.SandboxSpec{
					VMRuntime: "cloud-hypervisor",
					Image:     "ubuntu-24.04",
				},
			},
			want: false,
		},
		{
			name:     "image on another host",
			hostName: "host-b",
			sandbox: infrav1.Sandbox{
				Spec: infrav1.SandboxSpec{
					VMRuntime: "tart",
					Image:     "macos-14-base",
				},
			},
			want: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := hostHasReadyImageForSandbox(testCase.hostName, testCase.sandbox, hostImages)
			if got != testCase.want {
				t.Fatalf("hostHasReadyImageForSandbox() = %t, want %t", got, testCase.want)
			}
		})
	}
}
