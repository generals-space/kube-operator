package pod

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

func onAdd(obj interface{}) {
	pod := obj.(*corev1.Pod)
	klog.Infof("add a pod:%+v", pod.Name)
}

// onUpdate // 此处省略 workqueue 的使用
func onUpdate(oldObj interface{}, newObj interface{}) {
	klog.Infof("update a pod")
	klog.Infof("old object: %+v\n", oldObj)
	klog.Infof("new object: %+v\n", newObj)
}

func onDelete(obj interface{}) {
	klog.Infof("delete a pod")
}

type PodController struct {
	kubeClient *kubernetes.Clientset
}

func New(kubeClient *kubernetes.Clientset) *PodController {
	return &PodController{
		kubeClient: kubeClient,
	}
}

func (c *PodController) Start(stopCh chan struct{}) {
	lw := &cache.ListWatch{}
	resyncPeriod := time.Second * 60
	pod := &corev1.Pod{}
	podInformer := cache.NewSharedIndexInformer(lw, pod, resyncPeriod, cache.Indexers{
		cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
	})
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
		errTimeout := fmt.Errorf("初次同步缓存超时失败")
		runtime.HandleError(errTimeout)
		return
	}
	podInformer.Run(stopCh)

	// 这种方式创建的 informer 只能下面的语句得到 Lister 对象.
	podLister := listerv1.NewNodeLister(podInformer.GetIndexer())
	// 从 lister 中获取所有 items
	podList, err := podLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("获取Pod列表失败: %s", err)
	}
	klog.Infof("获取Pod列表: %+v", podList)
}
