/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlCli "sigs.k8s.io/controller-runtime/pkg/client"

	kubegroupv1 "generals-space/kube-operator/api/v1"
)

var logger = ctrl.Log.WithName("podcluster")

// PodClusterReconciler reconciles a PodCluster object
type PodClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kubegroup.generals.space,resources=podclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubegroup.generals.space,resources=podclusters/status,verbs=get;update;patch

func (r *PodClusterReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	result = ctrl.Result{}
	ctx := context.Background()

	// 这里不再展示 PodCluster 的相关流程, 只演示索引器的使用方法.

	// 假设存在某个节点, 名称为 nodeName.
	nodeName := "k8s-master-01"
	// 查询该节点上的所有 Pod.
	// listOpts := &ctrlCli.ListOptions{
	// 	FieldSelector: fields.SelectorFromSet(fields.Set{
	// 		".spec.nodeName": nodeName,
	// 	}),
	// 	Namespace: "",
	// }
	listOpts := ctrlCli.MatchingFields{
		// 此处的字段路径需要与索引器的路径保持一致.
		".spec.nodeName": nodeName,
	}
	// 注意: 如果 listOpts 中不包含 Namespace 条件, 则会在所有命名空间中查询.
	podList := &apicorev1.PodList{}
	// List 方法的实现在
	// controller-runtime/pkg/cache/internal/cache_reader.go -> CacheReader.List()
	// 其中会根据 FieldSelector 参数是否存在决定是否从索引器中进行查询.
	err = r.Client.List(ctx, podList, listOpts)
	if err != nil {
		if errors.IsNotFound(err) {
			// 可能是删除操作
			logger.Error(err, "pod list not found", "info", podList)
		}
		return result, err
	}

	return result, nil
}

func (r *PodClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// 以 .spec.nodename 字段路径为 key 为 Pod 创建索引, 之后在 controller 运行期间,
	// 可以使用 List() 方法配合 fieldSelector 参数查询指定 node 上的所有 Pod 的列表.
	//
	// kubebuilder 默认只加载了 namespace 索引器, 对应 ctrlCli.ListOptions.Namespace 的查询条件.
	if err := mgr.GetFieldIndexer().IndexField(
		&apicorev1.Pod{}, ".spec.nodeName",
		func(rawObj runtime.Object) []string {
			pod := rawObj.(*apicorev1.Pod)
			return []string{string(pod.Spec.NodeName)}
		},
	); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&kubegroupv1.PodCluster{}).
		Complete(r)
}
