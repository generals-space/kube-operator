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
	"fmt"
	"strconv"

	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apimmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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

func (r *PodClusterReconciler) reconcilePodForCluster(
	ctx context.Context,
	podcluster *kubegroupv1.PodCluster,
	currentPodList *apicorev1.PodList,
) (err error) {
	name, ns := podcluster.Name, podcluster.Namespace
	podReplicas := podcluster.Spec.PodReplicas
	logger.Info("create pod for podcluster", "namespace", ns, "podcluster", name)

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
	delOpts := &ctrlCli.DeleteOptions{}
	err = r.Client.Delete(ctx, currentPodList, delOpts)
	for i := 0; i < int(podReplicas); i++ {
		pod := podTemplate.DeepCopy()
		pod.ObjectMeta.Name = fmt.Sprintf("%s-%d", name, i)

		// SetControllerReference 只修改但不更新, 需要手动更新绑定关系才行,
		// 这里放在 Create() 之前, 如果是先 Create, 再想添加绑定关系, 则需要先 Get() 出来才行.
		err = ctrl.SetControllerReference(podcluster, pod, r.Scheme)
		if err != nil {
			logger.Error(err, "set controller reference failed", "podcluster", podcluster)
			return err
		}

		err = r.Create(ctx, pod)
		if err != nil {
			logger.Error(err, "create pod failed", "podcluster", podcluster)
			return err
		}
	}
	return
}

// +kubebuilder:rbac:groups=kubegroup.generals.space,resources=podclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubegroup.generals.space,resources=podclusters/status,verbs=get;update;patch

func (r *PodClusterReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	result = ctrl.Result{}
	ctx := context.Background()

	// 首先获取目标对象
	podcluster := &kubegroupv1.PodCluster{}
	// req.NamespacedName 的格式为 kube-system/podcluster-sample,
	// 配合第3个参数确定类型, 能够唯一确定目标对象
	// 可以通过 types.NamespacedName{
	// 	Name:      podcluster.Name,
	// 	Namespace: podcluster.Namespace,
	// } 声明
	err = r.Get(ctx, req.NamespacedName, podcluster)
	if err != nil {
		if errors.IsNotFound(err) {
			// 可能是删除操作
			logger.Info("podcluster object not found, maybe deleted", "info", podcluster)
			return result, nil
		}
		return result, err
	}
	logger.Info("podcluster object found", "info", podcluster)

	////////////////////////////////////////////////////////////////////////////
	// 如果目标对象正在被删除(貌似 Pod 的 Terminating 状态就是这种情况)
	if podcluster.DeletionTimestamp != nil {
		logger.Info("get deleted podcluster, do nothing", "info", podcluster)
		return result, nil
	}

	// 获取属于当前 podcluster 对象的 pod 列表, 之后调用 reconcilePodForCluster() 方法,
	// 让 pod 数量与 podcluster 中预期的数量保持一致
	listOpts := &ctrlCli.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			"podcluster": podcluster.Name,
		}),
	}
	podList := &apicorev1.PodList{}
	err = r.Client.List(ctx, podList, listOpts)
	if err != nil {
		if errors.IsNotFound(err) {
			// 可能是删除操作
			logger.Info("pods of cluster not found, create", "info", podcluster)
		}
		return result, err
	}
	podNum := strconv.Itoa(podList.Size())
	logger.Info("found pods of cluster update", "info", podNum)

	err = r.reconcilePodForCluster(ctx, podcluster, podList)
	if err != nil {
		logger.Error(err, "reconcile pod for cluster", podcluster.Name)
	}
	return result, nil
}

func (r *PodClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubegroupv1.PodCluster{}).
		Complete(r)
}
