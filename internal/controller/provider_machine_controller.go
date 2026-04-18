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

	log.Info("provider machine reconciliation placeholder", "providerKind", machine.Spec.Provider.Kind)
	return ctrl.Result{}, nil
}

func (r *ProviderMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.ProviderMachine{}).
		Named("provider-machine").
		Complete(r)
}
