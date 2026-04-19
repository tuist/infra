package controller

import (
	"context"

	infrav1 "github.com/tuist/infra/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type HostReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *HostReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx).WithValues("host", req.NamespacedName)

	var host infrav1.Host
	if err := r.Get(ctx, req.NamespacedName, &host); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var hostImages infrav1.HostImageList
	if err := r.List(ctx, &hostImages); err != nil {
		return ctrl.Result{}, err
	}

	readyImages := int32(0)
	pullingImages := int32(0)
	failedImages := int32(0)

	for _, hostImage := range hostImages.Items {
		if hostImage.Spec.HostRef != host.Name {
			continue
		}

		switch hostImage.Status.Phase {
		case infrav1.PhaseReady:
			readyImages++
		case infrav1.PhasePulling, infrav1.PhasePending:
			pullingImages++
		case infrav1.PhaseFailed:
			failedImages++
		}
	}

	updated := host
	updated.Status.ReadyImageCount = readyImages
	updated.Status.PullingImageCount = pullingImages
	updated.Status.FailedImageCount = failedImages

	if !host.Status.Healthy && host.Status.Phase == "" {
		updated.Status.Phase = infrav1.PhaseRegistered
	}

	if updated.Status.LastHeartbeat == nil {
		updated.Status.LastHeartbeat = now()
	}

	if updated.Status.Phase == infrav1.PhaseRegistered && updated.Status.Healthy {
		updated.Status.Phase = infrav1.PhaseReady
	}

	if updated.Status.Phase == infrav1.PhaseReady && !updated.Status.Healthy {
		updated.Status.Phase = infrav1.PhaseRegistered
	}

	if !hostStatusEqual(host.Status, updated.Status) {
		if err := r.Status().Update(ctx, &updated); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info(
		"host reconciliation complete",
		"readyImageCount", updated.Status.ReadyImageCount,
		"pullingImageCount", updated.Status.PullingImageCount,
		"failedImageCount", updated.Status.FailedImageCount,
		"phase", updated.Status.Phase,
	)

	return ctrl.Result{}, nil
}

func (r *HostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.Host{}).
		Named("host").
		Complete(r)
}
