package svc

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

func onAdd(obj interface{}) {
	svc := obj.(*corev1.Service)
	klog.Infof("add a service: %+v", svc.Name)
}

// onUpdate // 此处省略 workqueue 的使用
func onUpdate(oldObj interface{}, newObj interface{}) {
	klog.Infof("update a service")
	klog.Infof("old object: %+v\n", oldObj)
	klog.Infof("new object: %+v\n", newObj)
}

func onDelete(obj interface{}) {
	klog.Infof("delete a service")
}

type SvcController struct {
	kubeClient *kubernetes.Clientset
}

func New(kubeClient *kubernetes.Clientset) *SvcController {
	return &SvcController{
		kubeClient: kubeClient,
	}
}

func (c *SvcController) Start(stopCh chan struct{}) {
	klog.Infof("初始化 service informer...")
	// Shared指的是多个 lister 共享同一个cache, 而且资源的变化会同时通知到cache和listers.

	resync := time.Second * 60
	// 如下是为监听资源添加过滤选项的3种方法(第1, 2种已经注释掉了)
	// 第1种
	// factory := informers.NewSharedInformerFactory(c.kubeClient, resync)

	// listWatchOption 最终会被加工成 List() 或 Watch() 方法可接受的参数类型, 
	// 如 client.CoreV1().Pods(namespace).List(options)
	listWatchOption := func(opt *metav1.ListOptions) {
		labelSet := labels.Set(map[string]string{
			"k8s-app": "kube-dns",
		})
		opt.LabelSelector = labels.SelectorFromSet(labelSet).String()
	}
	// 第2种
	// filterOpts := informers.WithTweakListOptions(listWatchOption)
	// factory := informers.NewSharedInformerFactoryWithOptions(c.kubeClient, resync, filterOpts)
	
	// 第3种
	factory := informers.NewFilteredSharedInformerFactory(c.kubeClient, resync, "kube-system", listWatchOption)

	// factory 初始时是没有监听任何资源的, 只有在获取某一种资源的 Informer 时(比如 Node),
	// 才会调用 factory.InformerFor(&corev1.Service{}) 将其加入到监听列表,
	// 然后调用 Start() 和 WaitForCacheSync() 才有意义.
	// 且这些各种资源的 Informer 其实都是借助 sharedIndexInformer 实现的.
	// cmInformer 拥有两个方法: Informer, Lister. 而 Informer 其实就是 watch 操作.
	cmInformer := factory.Core().V1().Services()
	informer := cmInformer.Informer()
	defer runtime.HandleCrash()

	// 启动 informer, 开始 list & watch 流程
	go factory.Start(stopCh)

	// 从 apiserver 同步某种资源的全部对象, 即 list.
	// 之后就可以使用watch这种资源, 维护这份缓存.
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		errTimeout := fmt.Errorf("初次同步缓存超时失败")
		runtime.HandleError(errTimeout)
		return
	}

	// 使用自定义 handler, 处理 watch 响应的各种事件.
	// 具体的维护操作在informer内部执行, 这里挂载的是额外的触发器.
	// 需要注意的是, 在上面的list过程中, 会不断触发onAdd事件, 相当于服务发现了.
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	// 从informer对象创建lister, 不过这里的代码没有特殊的目的,
	// 应该只是展示一下通过informer的接口得到list资源的方法.
	svcLister := cmInformer.Lister()
	// 从 lister 中获取所有 items
	svcList, err := svcLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("获取主机列表失败: %s", err)
	}
	klog.Infof("获取主机列表: %+v", svcList)
}
