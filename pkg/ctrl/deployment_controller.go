package ctrl

import (
	context "context"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info().Msgf("Reconciling Deployment: %s/%s", req.Namespace, req.Name)
	return ctrl.Result{}, nil
}

func AddDeploymentController(mgr manager.Manager) error {
	r := &DeploymentReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
