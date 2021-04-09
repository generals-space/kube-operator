package podcluster

import (
	"context"
	"fmt"

	kubegroupv1 "generals-space/kube-operator/pkg/apis/kubegroup/v1"

	"k8s.io/klog"
	apicorev1 "k8s.io/api/core/v1"
	apimmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"k8s.io/apimachinery/pkg/labels"
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

func (r *ReconcilePodCluster) reconcilePodForCluster(
	ctx context.Context,
	podcluster *kubegroupv1.PodCluster,
	currentPodList *apicorev1.PodList,
) (err error) {
	name, ns := podcluster.Name, podcluster.Namespace
	podReplicas := podcluster.Spec.PodReplicas
	klog.Infof("create pod for podcluster, namespace: %s, podcluster: %s", ns, name)

	// 如果 pod 数量与 podcluster 对象中的 PodReplicas 值一致, 则直接返回.
	if currentPodList.Size() == int(podReplicas) {
		return
	}

	podTemplate := &apicorev1.Pod{
		ObjectMeta: apimmetav1.ObjectMeta{
			Name:      "",
			Namespace: ns,
			Labels: map[string]string{
				"podcluster": name,
			},
		},
		Spec: apicorev1.PodSpec{
			Containers: []apicorev1.Container{
				{
					Name:    name,
					Image:   "registry.cn-hangzhou.aliyuncs.com/generals-space/centos7",
					Command: []string{"tail", "-f", "/etc/os-release"},
					Env: []apicorev1.EnvVar{
						{
							Name:  "PodReplicas",
							Value: string(podReplicas),
						},
					},
				},
			},
		},
	}

	// 如果 pod 数量与 podcluster 对象中的 PodReplicas 值一致, 则删除所有已有 pod, 然后重建.
	// 其实本来应该像sts那样多退少补的, 不过写起来太麻烦, 这里简化了这个过程
	delOpts := &client.DeleteOptions{}
	err = r.Client.Delete(ctx, currentPodList, delOpts)
	for i := 0; i < int(podReplicas); i++ {
		pod := podTemplate.DeepCopy()
		pod.ObjectMeta.Name = fmt.Sprintf("%s-%d", name, i)

		// SetControllerReference 只修改但不更新, 需要手动更新绑定关系才行,
		// 这里放在 Create() 之前, 如果是先 Create, 再想添加绑定关系, 则需要先 Get() 出来才行.
		err = ctrl.SetControllerReference(podcluster, pod, r.scheme)
		if err != nil {
			klog.Errorf("set controller reference failed for podcluster: %s failed", name)
			return err
		}

		err = r.Create(ctx, pod)
		if err != nil {
			klog.Errorf("create pod for podcluster: %s failed", name)
			return err
		}
	}
	return
}

// Reconcile reads that state of the cluster for a PodCluster object and makes changes based on the state read
// and what is in the PodCluster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// +kubebuilder:rbac:groups=kubegroup.generals.space,resources=podclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubegroup.generals.space,resources=podclusters/status,verbs=get;update;patch
func (r *ReconcilePodCluster) Reconcile(request reconcile.Request) (result reconcile.Result, err error) {
	result = ctrl.Result{}

	// Fetch the PodCluster instance
	// 首先获取目标对象
	podcluster := &kubegroupv1.PodCluster{}
	err = r.Get(context.TODO(), request.NamespacedName, podcluster)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return result, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	////////////////////////////////////////////////////////////////////////////
	// 如果目标对象正在被删除(貌似 Pod 的 Terminating 状态就是这种情况)
	if podcluster.DeletionTimestamp != nil {
		klog.Infof("get deleted podcluster %s, do nothing", podcluster.Name)
		return result, nil
	}

	// 获取属于当前 podcluster 对象的 pod 列表, 之后调用 reconcilePodForCluster() 方法,
	// 让 pod 数量与 podcluster 中预期的数量保持一致
	listOpts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			"podcluster": podcluster.Name,
		}),
	}
	podList := &apicorev1.PodList{}
	err = r.Client.List(context.Background(), podList, listOpts)
	if err != nil {
		if errors.IsNotFound(err) {
			// 可能是删除操作
			klog.Infof("pods of cluster %s not found, create", podcluster.Name)
		}
		return result, err
	}
	klog.Infof("found %d pods of cluster %s, update", podList.Size(), podcluster.Name)

	err = r.reconcilePodForCluster(context.Background(), podcluster, podList)
	if err != nil {
		klog.Errorf("reconcile pod for cluster %s failed: %s", podcluster.Name, err)
	}
	return result, nil
}
