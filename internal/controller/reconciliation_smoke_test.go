package controller

import (
	"context"
	"testing"

	infrav1 "github.com/tuist/infra/api/v1"
	infrascheme "github.com/tuist/infra/internal/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestReconciliationSmokeFlow(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := infrascheme.AddToScheme(scheme); err != nil {
		t.Fatalf("AddToScheme() returned error: %v", err)
	}

	pool := &infrav1.HostPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pool-a",
			Namespace: "default",
		},
		Spec: infrav1.HostPoolSpec{
			VMRuntimes: []string{"local"},
			WarmImages: []string{"ubuntu-24.04"},
			WarmPool: infrav1.WarmPoolSpec{
				MinIdleHosts: 1,
			},
		},
	}

	host := &infrav1.Host{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host-a",
			Namespace: "default",
		},
		Spec: infrav1.HostSpec{
			PoolRef: "pool-a",
			VMRuntimes: []infrav1.VMRuntimeCapability{
				{
					Name:      "local",
					GuestOSes: []string{"linux"},
				},
			},
			OS: "local",
		},
		Status: infrav1.HostStatus{
			Healthy:          true,
			Phase:            infrav1.PhaseRegistered,
			AllocatableSlots: 1,
		},
	}

	sandbox := &infrav1.Sandbox{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sandbox-a",
			Namespace: "default",
		},
		Spec: infrav1.SandboxSpec{
			VMRuntime: "local",
			GuestOS:   "linux",
			Image:     "ubuntu-24.04",
		},
	}

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(
			&infrav1.HostPool{},
			&infrav1.ProviderMachine{},
			&infrav1.Host{},
			&infrav1.HostImage{},
			&infrav1.Sandbox{},
		).
		WithObjects(pool, host, sandbox).
		Build()

	ctx := context.Background()

	hostPoolReconciler := &HostPoolReconciler{Client: cl, Scheme: scheme}
	if _, err := hostPoolReconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{Name: pool.Name, Namespace: pool.Namespace}}); err != nil {
		t.Fatalf("HostPool reconcile returned error: %v", err)
	}

	var updatedPool infrav1.HostPool
	if err := cl.Get(ctx, client.ObjectKeyFromObject(pool), &updatedPool); err != nil {
		t.Fatalf("Get(HostPool) returned error: %v", err)
	}

	if updatedPool.Status.Phase != infrav1.PhaseReady {
		t.Fatalf("HostPool phase = %q, want %q", updatedPool.Status.Phase, infrav1.PhaseReady)
	}

	if updatedPool.Status.DesiredHosts != 1 {
		t.Fatalf("HostPool desiredHosts = %d, want 1", updatedPool.Status.DesiredHosts)
	}

	var hostImages infrav1.HostImageList
	if err := cl.List(ctx, &hostImages, client.InNamespace("default")); err != nil {
		t.Fatalf("List(HostImage) returned error: %v", err)
	}

	if len(hostImages.Items) != 1 {
		t.Fatalf("HostImage count = %d, want 1", len(hostImages.Items))
	}

	hostImage := hostImages.Items[0]
	hostImageReconciler := &HostImageReconciler{Client: cl, Scheme: scheme}
	if _, err := hostImageReconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(&hostImage)}); err != nil {
		t.Fatalf("HostImage reconcile returned error: %v", err)
	}

	if err := cl.Get(ctx, client.ObjectKeyFromObject(&hostImage), &hostImage); err != nil {
		t.Fatalf("Get(HostImage) returned error: %v", err)
	}

	if hostImage.Status.Phase != infrav1.PhaseReady {
		t.Fatalf("HostImage phase = %q, want %q", hostImage.Status.Phase, infrav1.PhaseReady)
	}

	hostReconciler := &HostReconciler{Client: cl, Scheme: scheme}
	if _, err := hostReconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(host)}); err != nil {
		t.Fatalf("Host reconcile returned error: %v", err)
	}

	var updatedHost infrav1.Host
	if err := cl.Get(ctx, client.ObjectKeyFromObject(host), &updatedHost); err != nil {
		t.Fatalf("Get(Host) returned error: %v", err)
	}

	if updatedHost.Status.Phase != infrav1.PhaseReady {
		t.Fatalf("Host phase = %q, want %q", updatedHost.Status.Phase, infrav1.PhaseReady)
	}

	if updatedHost.Status.ReadyImageCount != 1 {
		t.Fatalf("Host readyImageCount = %d, want 1", updatedHost.Status.ReadyImageCount)
	}

	sandboxReconciler := &SandboxReconciler{Client: cl, Scheme: scheme}
	if _, err := sandboxReconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(sandbox)}); err != nil {
		t.Fatalf("Sandbox reconcile returned error: %v", err)
	}

	var updatedSandbox infrav1.Sandbox
	if err := cl.Get(ctx, client.ObjectKeyFromObject(sandbox), &updatedSandbox); err != nil {
		t.Fatalf("Get(Sandbox) returned error: %v", err)
	}

	if updatedSandbox.Status.Phase != infrav1.PhaseScheduled {
		t.Fatalf("Sandbox phase = %q, want %q", updatedSandbox.Status.Phase, infrav1.PhaseScheduled)
	}

	if updatedSandbox.Status.HostRef != "host-a" {
		t.Fatalf("Sandbox hostRef = %q, want %q", updatedSandbox.Status.HostRef, "host-a")
	}
}

func TestProviderMachineReconcileUpdatesObservedHost(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := infrascheme.AddToScheme(scheme); err != nil {
		t.Fatalf("AddToScheme() returned error: %v", err)
	}

	machine := &infrav1.ProviderMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "machine-a",
			Namespace: "default",
		},
		Spec: infrav1.ProviderMachineSpec{
			Provider: infrav1.ProviderRef{Kind: "aws-ec2-bare-metal"},
		},
		Status: infrav1.ProviderMachineStatus{
			ProviderID: "provider-123",
		},
	}

	host := &infrav1.Host{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host-a",
			Namespace: "default",
		},
		Spec: infrav1.HostSpec{
			ProviderMachineRef: "machine-a",
		},
	}

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(&infrav1.ProviderMachine{}).
		WithObjects(machine, host).
		Build()

	reconciler := &ProviderMachineReconciler{Client: cl, Scheme: scheme}
	if _, err := reconciler.Reconcile(context.Background(), ctrl.Request{NamespacedName: client.ObjectKeyFromObject(machine)}); err != nil {
		t.Fatalf("ProviderMachine reconcile returned error: %v", err)
	}

	var updatedMachine infrav1.ProviderMachine
	if err := cl.Get(context.Background(), client.ObjectKeyFromObject(machine), &updatedMachine); err != nil {
		t.Fatalf("Get(ProviderMachine) returned error: %v", err)
	}

	if updatedMachine.Status.Phase != infrav1.PhaseRegistered {
		t.Fatalf("ProviderMachine phase = %q, want %q", updatedMachine.Status.Phase, infrav1.PhaseRegistered)
	}

	if updatedMachine.Status.ObservedHost != "host-a" {
		t.Fatalf("ProviderMachine observedHost = %q, want %q", updatedMachine.Status.ObservedHost, "host-a")
	}
}
