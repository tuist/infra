package controller

import (
	"context"

	kubeboxv1 "github.com/tuist/kubebox/api/v1"
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

	var sandbox kubeboxv1.Sandbox
	if err := r.Get(ctx, req.NamespacedName, &sandbox); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var hosts kubeboxv1.HostList
	if err := r.List(ctx, &hosts); err != nil {
		return ctrl.Result{}, err
	}

	matchingHosts := make([]string, 0, len(hosts.Items))
	for _, host := range hosts.Items {
		if hostSupportsSandbox(host, sandbox) {
			matchingHosts = append(matchingHosts, host.Name)
		}
	}

	log.Info(
		"sandbox reconciliation placeholder",
		"tenant", sandbox.Spec.Tenant,
		"classRef", sandbox.Spec.ClassRef,
		"backend", sandbox.Spec.Backend,
		"guestOS", sandbox.Spec.GuestOS,
		"matchingHosts", matchingHosts,
	)
	return ctrl.Result{}, nil
}

func (r *SandboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeboxv1.Sandbox{}).
		Named("sandbox").
		Complete(r)
}
