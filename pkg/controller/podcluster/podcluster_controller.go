package podcluster

import (
	"context"

	kubegroupv1 "generals-space/kube-operator/pkg/apis/kubegroup/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new PodCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePodCluster{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("podcluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to PodCluster
	err = c.Watch(&source.Kind{Type: &kubegroupv1.PodCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by PodCluster - change this for objects you create
	// err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
	//	IsController: true,
	//	OwnerType:    &kubegroupv1.PodCluster{},
	// })
	// if err != nil {
	//	return err
	// }

	return nil
}

var _ reconcile.Reconciler = &ReconcilePodCluster{}

// ReconcilePodCluster reconciles a PodCluster object
type ReconcilePodCluster struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a PodCluster object and makes changes based on the state read
// and what is in the PodCluster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// +kubebuilder:rbac:groups=kubegroup.generals.space,resources=podclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubegroup.generals.space,resources=podclusters/status,verbs=get;update;patch
func (r *ReconcilePodCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the PodCluster instance
	instance := &kubegroupv1.PodCluster{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
