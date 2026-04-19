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

	var host infrav1.Host
	hostErr := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: hostImage.Spec.HostRef}, &host)
	if client.IgnoreNotFound(hostErr) != nil {
		return ctrl.Result{}, hostErr
	}

	updated := hostImage
	updated.Status.LastUpdateTime = now()

	switch {
	case hostImage.Spec.HostRef == "" || hostImage.Spec.VMRuntime == "" || hostImage.Spec.Image == "":
		updated.Status.Phase = infrav1.PhaseFailed
		updated.Status.FailureMessage = "hostRef, vmRuntime, and image are required"
		updated.Status.ProgressPercent = 0
	case hostErr != nil || !host.Status.Healthy:
		updated.Status.Phase = infrav1.PhasePending
		updated.Status.FailureMessage = ""
		updated.Status.ProgressPercent = 0
	case host.Spec.OS == "local":
		updated.Status.Phase = infrav1.PhaseReady
		updated.Status.FailureMessage = ""
		updated.Status.ProgressPercent = 100
	default:
		updated.Status.Phase = infrav1.PhasePulling
		updated.Status.FailureMessage = ""
		if updated.Status.ProgressPercent == 0 {
			updated.Status.ProgressPercent = 10
		}
	}

	if !hostImageStatusEqual(hostImage.Status, updated.Status) {
		if err := r.Status().Update(ctx, &updated); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info(
		"host image reconciliation complete",
		"hostRef", hostImage.Spec.HostRef,
		"vmRuntime", hostImage.Spec.VMRuntime,
		"image", hostImage.Spec.Image,
		"phase", updated.Status.Phase,
		"progressPercent", updated.Status.ProgressPercent,
	)
	return ctrl.Result{}, nil
}

func (r *HostImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.HostImage{}).
		Named("hostimage").
		Complete(r)
}
