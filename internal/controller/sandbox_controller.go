package controller

import (
	"context"

	kubeboxv1alpha1 "github.com/tuist/kubebox/api/v1alpha1"
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

	var sandbox kubeboxv1alpha1.Sandbox
	if err := r.Get(ctx, req.NamespacedName, &sandbox); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("sandbox reconciliation placeholder", "tenant", sandbox.Spec.Tenant, "classRef", sandbox.Spec.ClassRef)
	return ctrl.Result{}, nil
}

func (r *SandboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeboxv1alpha1.Sandbox{}).
		Named("sandbox").
		Complete(r)
}

