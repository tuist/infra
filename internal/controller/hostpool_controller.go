package controller

import (
	"context"

	kubeboxv1 "github.com/tuist/kubebox/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type HostPoolReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *HostPoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx).WithValues("hostpool", req.NamespacedName)

	var pool kubeboxv1.HostPool
	if err := r.Get(ctx, req.NamespacedName, &pool); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("host pool reconciliation placeholder", "desiredIdleHosts", pool.Spec.WarmPool.MinIdleHosts)
	return ctrl.Result{}, nil
}

func (r *HostPoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeboxv1.HostPool{}).
		Named("hostpool").
		Complete(r)
}
