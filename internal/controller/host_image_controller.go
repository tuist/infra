package controller

import (
	"context"

	infrav1 "github.com/tuist/infra/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type HostImageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *HostImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx).WithValues("hostImage", req.NamespacedName)

	var hostImage infrav1.HostImage
	if err := r.Get(ctx, req.NamespacedName, &hostImage); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info(
		"host image reconciliation placeholder",
		"hostRef", hostImage.Spec.HostRef,
		"vmRuntime", hostImage.Spec.VMRuntime,
		"image", hostImage.Spec.Image,
		"phase", hostImage.Status.Phase,
		"progressPercent", hostImage.Status.ProgressPercent,
	)
	return ctrl.Result{}, nil
}

func (r *HostImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.HostImage{}).
		Named("hostimage").
		Complete(r)
}
