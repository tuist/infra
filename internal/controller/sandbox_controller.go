package controller

import (
	"context"
	"slices"

	infrav1 "github.com/tuist/infra/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SandboxReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *SandboxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx).WithValues("sandbox", req.NamespacedName)

	var sandbox infrav1.Sandbox
	if err := r.Get(ctx, req.NamespacedName, &sandbox); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var hosts infrav1.HostList
	if err := r.List(ctx, &hosts); err != nil {
		return ctrl.Result{}, err
	}

	var hostImages infrav1.HostImageList
	if err := r.List(ctx, &hostImages); err != nil {
		return ctrl.Result{}, err
	}

	matchingHosts := make([]string, 0, len(hosts.Items))
	for _, host := range hosts.Items {
		if hostSupportsSandbox(host, sandbox) && hostHasReadyImageForSandbox(host.Name, sandbox, hostImages.Items) {
			matchingHosts = append(matchingHosts, host.Name)
		}
	}

	slices.Sort(matchingHosts)

	updated := sandbox
	if len(matchingHosts) == 0 {
		updated.Status.Phase = infrav1.PhasePending
		updated.Status.HostRef = ""
	} else {
		updated.Status.Phase = infrav1.PhaseScheduled
		updated.Status.HostRef = matchingHosts[0]
	}

	if !sandboxStatusEqual(sandbox.Status, updated.Status) {
		if err := r.Status().Update(ctx, &updated); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info(
		"sandbox reconciliation complete",
		"tenant", sandbox.Spec.Tenant,
		"classRef", sandbox.Spec.ClassRef,
		"vmRuntime", sandbox.Spec.VMRuntime,
		"image", sandbox.Spec.Image,
		"guestOS", sandbox.Spec.GuestOS,
		"selectedHost", updated.Status.HostRef,
		"phase", updated.Status.Phase,
		"matchingHosts", matchingHosts,
	)
	return ctrl.Result{}, nil
}

func (r *SandboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.Sandbox{}).
		Named("sandbox").
		Complete(r)
}
