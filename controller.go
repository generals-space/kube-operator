package main

import (
	"fmt"
	"time"

	"k8s.io/klog"
	apiCorev1 "k8s.io/api/core/v1"
	apimErrors "k8s.io/apimachinery/pkg/api/errors"
	apimRuntime "k8s.io/apimachinery/pkg/util/runtime"
	apimUtilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	apimWait "k8s.io/apimachinery/pkg/util/wait"
	cgKuber "k8s.io/client-go/kubernetes"
	cgScheme "k8s.io/client-go/kubernetes/scheme"
	cgCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	cgCache "k8s.io/client-go/tools/cache"
	cgRecord "k8s.io/client-go/tools/record"
	cgWorkqueue "k8s.io/client-go/util/workqueue"

	crdV1 "generals-space/kube-operator/pkg/apis/kubegroup/v1"
	crdClientset "generals-space/kube-operator/pkg/client/clientset/versioned"
	crdScheme "generals-space/kube-operator/pkg/client/clientset/versioned/scheme"
	crdInformers "generals-space/kube-operator/pkg/client/informers/externalversions/kubegroup/v1"
	crdListers "generals-space/kube-operator/pkg/client/listers/kubegroup/v1"
)

const controllerAgentName = "podcluster-controller"

const (
	// SuccessSynced ...
	SuccessSynced = "Synced"
	// MessageResourceSynced ...
	MessageResourceSynced = "PodCluster synced successfully"
)

// Controller is the controller implementation for PodCluster resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset cgKuber.Interface
	// podClusterClientset is a clientset for our own API group
	podClusterClientset crdClientset.Interface

	podClusterLister crdListers.PodClusterLister
	podClusterSynced cgCache.InformerSynced

	workqueue cgWorkqueue.RateLimitingInterface
	recorder  cgRecord.EventRecorder
}

// NewController returns a new pod group controller
func NewController(
	kubeclientset cgKuber.Interface,
	podClusterClientset crdClientset.Interface,
	podClusterInformer crdInformers.PodClusterInformer) *Controller {

	// AddToScheme() 将 CRD 的结构与 Kubernetes GroupVersionKinds 的映射添加到 Manager 的 Scheme 中
	// 以便能够让 Controller Manager 知道 CRD 的存在
	apimUtilRuntime.Must(crdScheme.AddToScheme(cgScheme.Scheme))
	klog.Info("Creating event broadcaster")
	eventBroadcaster := cgRecord.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(
		&cgCorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")},
	)
	recorder := eventBroadcaster.NewRecorder(
		cgScheme.Scheme, apiCorev1.EventSource{Component: controllerAgentName},
	)

	controller := &Controller{
		kubeclientset:       kubeclientset,
		podClusterClientset: podClusterClientset,
		podClusterLister:    podClusterInformer.Lister(),
		podClusterSynced:    podClusterInformer.Informer().HasSynced,
		workqueue:           cgWorkqueue.NewNamedRateLimitingQueue(
			cgWorkqueue.DefaultControllerRateLimiter(), "PodClusters",
		),
		recorder:            recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when PodCluster resources change
	podClusterInformer.Informer().AddEventHandler(cgCache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuePodCluster,
		UpdateFunc: func(old, new interface{}) {
			oldPodCluster := old.(*crdV1.PodCluster)
			newPodCluster := new.(*crdV1.PodCluster)
			if oldPodCluster.ResourceVersion == newPodCluster.ResourceVersion {
				//版本一致, 就表示没有实际更新的操作, 立即返回
				return
			}
			controller.enqueuePodCluster(new)
		},
		DeleteFunc: controller.enqueuePodClusterForDelete,
	})

	return controller
}

// Run 在此处开始controller的业务
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer apimRuntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("开始controller业务, 开始一次缓存数据同步")
	if ok := cgCache.WaitForCacheSync(stopCh, c.podClusterSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("worker启动")
	for i := 0; i < threadiness; i++ {
		go apimWait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("worker已经启动")
	<-stopCh
	klog.Info("worker已经结束")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// 取数据处理
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			apimRuntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// 在syncHandler中处理业务
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		apimRuntime.HandleError(err)
		return true
	}

	return true
}

// 处理
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cgCache.SplitMetaNamespaceKey(key)
	if err != nil {
		apimRuntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// 从缓存中取对象
	podCluster, err := c.podClusterLister.PodClusters(namespace).Get(name)
	if err != nil {
		// 如果PodCluster对象被删除了, 就会走到这里, 所以应该在这里加入执行
		if apimErrors.IsNotFound(err) {
			klog.Infof("PodCluster对象被删除, 请在这里执行实际的删除业务: %s/%s ...", namespace, name)
			return nil
		}
		apimRuntime.HandleError(fmt.Errorf("failed to list podCluster by: %s/%s", namespace, name))
		return err
	}

	klog.Infof("这里是podCluster对象的期望状态: %#v ...", podCluster)
	klog.Infof("实际状态是从业务层面得到的, 此处应该去的实际状态, 与期望状态做对比, 并根据差异做出响应(新增或者删除)")

	c.recorder.Event(podCluster, apiCorev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// 数据先放入缓存, 再入队列
func (c *Controller) enqueuePodCluster(obj interface{}) {
	var key string
	var err error
	// 将对象放入缓存
	if key, err = cgCache.MetaNamespaceKeyFunc(obj); err != nil {
		apimRuntime.HandleError(err)
		return
	}

	// 将key放入队列
	c.workqueue.AddRateLimited(key)
}

// 删除操作
func (c *Controller) enqueuePodClusterForDelete(obj interface{}) {
	var key string
	var err error
	// 从缓存中删除指定对象
	key, err = cgCache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		apimRuntime.HandleError(err)
		return
	}
	//再将key放入队列
	c.workqueue.AddRateLimited(key)
}
