package controller

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"slices"
	"time"

	infrav1 "github.com/tuist/infra/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	labelPoolRef   = "infra.tuist.dev/pool"
	labelHostRef   = "infra.tuist.dev/host"
	labelVMRuntime = "infra.tuist.dev/vm-runtime"
)

func now() *metav1.Time {
	timestamp := metav1.NewTime(time.Now().UTC())
	return &timestamp
}

func hostBelongsToPool(host infrav1.Host, pool infrav1.HostPool) bool {
	return host.Spec.PoolRef == pool.Name
}

func hostMatchesProviderMachine(host infrav1.Host, machine infrav1.ProviderMachine) bool {
	return host.Spec.ProviderMachineRef == machine.Name
}

func shouldCountHostReady(host infrav1.Host) bool {
	return host.Status.Healthy && (host.Status.Phase == "" || host.Status.Phase == infrav1.PhaseReady || host.Status.Phase == infrav1.PhaseRegistered)
}

func buildDesiredHostImages(pool infrav1.HostPool, hosts []infrav1.Host) map[string]infrav1.HostImage {
	desired := make(map[string]infrav1.HostImage)

	for _, host := range hosts {
		if !hostBelongsToPool(host, pool) {
			continue
		}

		for _, vmRuntime := range host.Spec.VMRuntimes {
			if !slices.Contains(pool.Spec.VMRuntimes, vmRuntime.Name) {
				continue
			}

			for _, image := range pool.Spec.WarmImages {
				hostImage := infrav1.HostImage{
					ObjectMeta: metav1.ObjectMeta{
						Name:      hostImageName(host.Name, vmRuntime.Name, image),
						Namespace: pool.Namespace,
						Labels: map[string]string{
							labelPoolRef:   pool.Name,
							labelHostRef:   host.Name,
							labelVMRuntime: vmRuntime.Name,
						},
					},
					Spec: infrav1.HostImageSpec{
						HostRef:   host.Name,
						VMRuntime: vmRuntime.Name,
						Image:     image,
					},
				}

				desired[hostImage.Name] = hostImage
			}
		}
	}

	return desired
}

func hostImageName(hostName, vmRuntime, image string) string {
	sum := sha1.Sum([]byte(fmt.Sprintf("%s|%s|%s", hostName, vmRuntime, image)))
	return fmt.Sprintf("%s-%s-%s", hostName, vmRuntime, hex.EncodeToString(sum[:4]))
}

func hostPoolStatusEqual(left, right infrav1.HostPoolStatus) bool {
	return left.Phase == right.Phase &&
		left.DesiredHosts == right.DesiredHosts &&
		left.ReadyHosts == right.ReadyHosts &&
		left.AvailableSlots == right.AvailableSlots
}

func providerMachineStatusEqual(left, right infrav1.ProviderMachineStatus) bool {
	return left.Phase == right.Phase &&
		left.ProviderID == right.ProviderID &&
		left.ObservedHost == right.ObservedHost
}

func hostStatusEqual(left, right infrav1.HostStatus) bool {
	return left.Phase == right.Phase &&
		left.Healthy == right.Healthy &&
		left.AllocatableSlots == right.AllocatableSlots &&
		left.ReadyImageCount == right.ReadyImageCount &&
		left.PullingImageCount == right.PullingImageCount &&
		left.FailedImageCount == right.FailedImageCount
}

func hostImageStatusEqual(left, right infrav1.HostImageStatus) bool {
	return left.Phase == right.Phase &&
		left.ProgressPercent == right.ProgressPercent &&
		left.Digest == right.Digest &&
		left.FailureMessage == right.FailureMessage
}

func sandboxStatusEqual(left, right infrav1.SandboxStatus) bool {
	return left.Phase == right.Phase &&
		left.HostRef == right.HostRef &&
		left.LeaseRef == right.LeaseRef &&
		left.AccessURL == right.AccessURL
}
