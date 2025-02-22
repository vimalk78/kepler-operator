package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/sustainable.computing.io/kepler-operator/pkg/api/v1alpha1"
	"github.com/sustainable.computing.io/kepler-operator/pkg/components"
	"github.com/sustainable.computing.io/kepler-operator/pkg/components/exporter"
	"github.com/sustainable.computing.io/kepler-operator/pkg/reconciler"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"

	secv1 "github.com/openshift/api/security/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	KeplerFinalizer = "kepler.system.sustainable.computing.io/finalizer"
)

// KeplerReconciler reconciles a Kepler object
type KeplerReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	logger logr.Logger
}

// Owned resource
//+kubebuilder:rbac:groups=kepler.system.sustainable.computing.io,resources=*,verbs=*

// common to all components deployed by operator
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services;configmaps;serviceaccounts,verbs=list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=*,verbs=*

// RBAC for running Kepler exporter
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=security.openshift.io,resources=securitycontextconstraints,verbs=list;watch;create;update;patch;delete;use

// RBAC required by Kepler exporter
//+kubebuilder:rbac:groups=core,resources=nodes/metrics;nodes/proxy;nodes/stats,verbs=get;list;watch

// SetupWithManager sets up the controller with the Manager.
func (r *KeplerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// We only want to trigger a reconciliation when the generation
	// of a child changes. Until we need to update our the status for our own objects,
	// we can save CPU cycles by avoiding reconciliations triggered by
	// child status changes.
	//
	// TODO: consider using ResourceVersionChanged predicate for resources that support it

	genChanged := builder.WithPredicates(predicate.GenerationChangedPredicate{})

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Kepler{}).
		Owns(&corev1.ConfigMap{}, genChanged).
		Owns(&corev1.ServiceAccount{}, genChanged).
		Owns(&corev1.Service{}, genChanged).
		Owns(&appsv1.DaemonSet{}, genChanged).
		Owns(&rbacv1.ClusterRoleBinding{}, genChanged).
		Owns(&rbacv1.ClusterRole{}, genChanged).
		Owns(&secv1.SecurityContextConstraints{}, genChanged).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kepler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KeplerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO: remove these keys from the log
	// "controller": "kepler", "controllerGroup": "kepler.system.sustainable.computing.io",
	// "controllerKind": "Kepler", "Kepler": {"name":"kepler"},

	logger := log.FromContext(ctx)
	r.logger = logger

	logger.Info("Start of  reconcile")
	defer logger.Info("End of reconcile")

	kepler, err := r.getKepler(ctx, req)
	if err != nil {
		// retry since some error has occurred
		logger.V(6).Info("Get Error ", "error", err)
		return ctrl.Result{}, err
	}

	if kepler == nil {
		// no kepler found , so stop here
		logger.V(6).Info("Kepler Nil")
		return ctrl.Result{}, nil
	}

	// TODO: have admission webhook decline all but `kepler` instance
	if kepler.Name != "kepler" {
		return r.setInvalidStatus(ctx, req)
	}

	logger.V(6).Info("Running sub reconcilers", "kepler", kepler.Spec)

	result, recErr := r.runKeplerReconcilers(ctx, kepler)
	updateErr := r.updateStatus(ctx, req, err)

	if recErr != nil {
		return result, recErr
	}
	return result, updateErr
}

func (r KeplerReconciler) runKeplerReconcilers(ctx context.Context, kepler *v1alpha1.Kepler) (ctrl.Result, error) {

	reconcilers := r.reconcilersForKepler(kepler)
	r.logger.V(6).Info("renconcilers ...", "count", len(reconcilers))

	return reconciler.Runner{
		Reconcilers: reconcilers,
		Client:      r.Client,
		Scheme:      r.Scheme,
		Logger:      r.logger,
	}.Run(ctx)
}

func (r KeplerReconciler) getKepler(ctx context.Context, req ctrl.Request) (*v1alpha1.Kepler, error) {
	logger := r.logger

	kepler := v1alpha1.Kepler{}

	if err := r.Client.Get(ctx, req.NamespacedName, &kepler); err != nil {
		if errors.IsNotFound(err) {
			logger.V(3).Info("kepler could not be found; may be marked for deletion")
			return nil, nil
		}
		logger.Error(err, "failed to get kepler")
		return nil, err
	}

	return &kepler, nil
}

func (r KeplerReconciler) updateStatus(ctx context.Context, req ctrl.Request, recErr error) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {

		kepler, _ := r.getKepler(ctx, req)
		// may be deleted
		if kepler == nil || !kepler.GetDeletionTimestamp().IsZero() {
			// retry since some error has occurred
			r.logger.V(6).Info("Reconcile has deleted kepler; skipping update")
			return nil
		}

		if recErr != nil {
			kepler.Status.Conditions = []metav1.Condition{{
				Type:               v1alpha1.ConditionReconciled,
				Status:             metav1.ConditionTrue,
				ObservedGeneration: kepler.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             v1alpha1.ReconciledReasonError,
				Message:            recErr.Error(),
			}}
		} else {
			kepler.Status.Conditions = []metav1.Condition{{
				Type:               v1alpha1.ConditionReconciled,
				Status:             metav1.ConditionTrue,
				ObservedGeneration: kepler.Generation,
				Reason:             v1alpha1.ReconciledReasonComplete,
				Message:            "Reconcile succeeded",
				LastTransitionTime: metav1.Now(),
			}}
		}
		return r.Client.Status().Update(ctx, kepler)

	})
}

func (r KeplerReconciler) reconcilersForKepler(k *v1alpha1.Kepler) []reconciler.Reconciler {
	// start with namespace
	rs := []reconciler.Reconciler{}

	if cleanup := !k.DeletionTimestamp.IsZero(); cleanup {
		rs = append(rs, deletersForResources(components.NewKeplerNamespace())...)
	} else {
		rs = append(rs, reconciler.Updater{Owner: k, Resource: components.NewKeplerNamespace()})
	}

	rs = append(rs, exporterReconcilers(k)...)

	// TODO: add this when modelServer is supported by Kepler Spec
	// rs = append(rs, modelServerReconcilers(k)...)

	// Add/Remove finalizer at the end
	rs = append(rs, reconciler.Finalizer{Resource: k, Finalizer: KeplerFinalizer, Logger: r.logger})
	return rs
}

func (r KeplerReconciler) setInvalidStatus(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		invalidKepler, _ := r.getKepler(ctx, req)
		// may be deleted
		if invalidKepler == nil || !invalidKepler.GetDeletionTimestamp().IsZero() {
			return nil
		}

		invalidKepler.Status.Conditions = []metav1.Condition{{
			Type:               v1alpha1.ConditionReconciled,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: invalidKepler.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             v1alpha1.InvalidKeplerObjectReason,
			Message:            "Only a single instance of Kepler named kepler is allowed",
		}}
		return r.Client.Status().Update(ctx, invalidKepler)
	})

	// retry only on error
	return ctrl.Result{}, err
}

func exporterReconcilers(k *v1alpha1.Kepler) []reconciler.Reconciler {

	// TODO: validate if this required anymore after the cluster-scoped CRD change
	if cleanup := !k.DeletionTimestamp.IsZero(); cleanup {
		return deletersForResources(
			exporter.NewSCC(components.Metadata, k),
			exporter.NewClusterRoleBinding(components.Metadata),
			exporter.NewClusterRole(components.Metadata),
		)
	}

	return updatersForResources(k,
		// cluster-scoped resources first
		exporter.NewSCC(components.Full, k),
		exporter.NewClusterRole(components.Full),
		exporter.NewClusterRoleBinding(components.Full),

		// namespace scoped
		exporter.NewServiceAccount(),
		exporter.NewConfigMap(components.Full, k),
		exporter.NewDaemonSet(k),
		exporter.NewService(k),
	)
}

func updatersForResources(k *v1alpha1.Kepler, resources ...client.Object) []reconciler.Reconciler {
	rs := []reconciler.Reconciler{}
	resourceUpdater := newUpdaterForKepler(k)
	for _, res := range resources {
		rs = append(rs, resourceUpdater(res))
	}
	return rs

}

func deletersForResources(resources ...client.Object) []reconciler.Reconciler {
	rs := []reconciler.Reconciler{}
	for _, res := range resources {
		rs = append(rs, reconciler.Deleter{Resource: res})
	}
	return rs
}

// TODO: decide if this this should move to reconciler
type newUpdaterFn func(client.Object) *reconciler.Updater

func newUpdaterForKepler(k *v1alpha1.Kepler) newUpdaterFn {
	return func(obj client.Object) *reconciler.Updater {
		// NOTE: Owner: k.GetObjectMeta() also works
		return &reconciler.Updater{Owner: k, Resource: obj}
	}
}
