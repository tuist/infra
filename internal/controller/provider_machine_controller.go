package controller

import (
	"context"

	infrav1 "github.com/tuist/infra/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ProviderMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ProviderMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx).WithValues("providerMachine", req.NamespacedName)

	var machine infrav1.ProviderMachine
	if err := r.Get(ctx, req.NamespacedName, &machine); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var hosts infrav1.HostList
	if err := r.List(ctx, &hosts, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	updated := machine
	updated.Status.LastUpdateTime = now()

	foundHost := false
	for _, host := range hosts.Items {
		if !hostMatchesProviderMachine(host, machine) {
			continue
		}

		updated.Status.ObservedHost = host.Name
		updated.Status.Phase = infrav1.PhaseRegistered
		foundHost = true
		break
	}

	if !foundHost {
		updated.Status.ObservedHost = ""
		if machine.Status.ProviderID != "" {
			updated.Status.Phase = infrav1.PhaseBootstrapping
		} else {
			updated.Status.Phase = infrav1.PhasePending
		}
	}

	if !providerMachineStatusEqual(machine.Status, updated.Status) {
		if err := r.Status().Update(ctx, &updated); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info(
		"provider machine reconciliation complete",
		"providerKind", machine.Spec.Provider.Kind,
		"phase", updated.Status.Phase,
		"observedHost", updated.Status.ObservedHost,
	)
	return ctrl.Result{}, nil
}

func (r *ProviderMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.ProviderMachine{}).
		Named("provider-machine").
		Complete(r)
}
