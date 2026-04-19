package controller

import (
	"context"

	infrav1 "github.com/tuist/infra/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type HostPoolReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *HostPoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx).WithValues("hostpool", req.NamespacedName)

	var pool infrav1.HostPool
	if err := r.Get(ctx, req.NamespacedName, &pool); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var hosts infrav1.HostList
	if err := r.List(ctx, &hosts, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	var hostImages infrav1.HostImageList
	if err := r.List(ctx, &hostImages, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	desiredHostImages := buildDesiredHostImages(pool, hosts.Items)
	existingHostImages := make(map[string]infrav1.HostImage)
	for _, hostImage := range hostImages.Items {
		if hostImage.Labels[labelPoolRef] != pool.Name {
			continue
		}

		existingHostImages[hostImage.Name] = hostImage
	}

	for name, desiredHostImage := range desiredHostImages {
		if _, exists := existingHostImages[name]; exists {
			continue
		}

		if err := controllerutil.SetControllerReference(&pool, &desiredHostImage, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, &desiredHostImage); err != nil {
			return ctrl.Result{}, err
		}
	}

	for name, existingHostImage := range existingHostImages {
		if _, exists := desiredHostImages[name]; exists {
			continue
		}

		if err := r.Delete(ctx, &existingHostImage); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
	}

	desiredHosts := pool.Spec.WarmPool.MinIdleHosts
	readyHosts := int32(0)
	availableSlots := int32(0)
	for _, host := range hosts.Items {
		if !hostBelongsToPool(host, pool) {
			continue
		}

		if shouldCountHostReady(host) {
			readyHosts++
		}

		availableSlots += host.Status.AllocatableSlots
	}

	updated := pool
	updated.Status.DesiredHosts = desiredHosts
	updated.Status.ReadyHosts = readyHosts
	updated.Status.AvailableSlots = availableSlots
	updated.Status.LastScaleTime = now()

	switch {
	case desiredHosts == 0:
		updated.Status.Phase = infrav1.PhaseReady
	case readyHosts >= desiredHosts:
		updated.Status.Phase = infrav1.PhaseReady
	case readyHosts > 0:
		updated.Status.Phase = infrav1.PhaseProvisioning
	default:
		updated.Status.Phase = infrav1.PhasePending
	}

	if !hostPoolStatusEqual(pool.Status, updated.Status) {
		if err := r.Status().Update(ctx, &updated); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info(
		"host pool reconciliation complete",
		"desiredHosts", desiredHosts,
		"readyHosts", readyHosts,
		"availableSlots", availableSlots,
		"managedHostImages", len(desiredHostImages),
		"phase", updated.Status.Phase,
	)
	return ctrl.Result{}, nil
}

func (r *HostPoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.HostPool{}).
		Named("hostpool").
		Complete(r)
}
