package main

import (
	"log"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog"
)

////////////////////////////////////////////////////////////////////////////////
// 各种处理函数中, obj即为watch接口响应得到的资源对象.

func onAdd(obj interface{}) {
	pod := obj.(*corev1.Pod)
	klog.Infof("add a pod: %+v", pod.Name)
}

// onUpdate // 此处省略 workqueue 的使用
func onUpdate(oldObj interface{}, newObj interface{}) {
	klog.Infof("update a pod")
	oldPod := oldObj.(*corev1.Pod)
	newPod := newObj.(*corev1.Pod)

	klog.Infof("old pod: %+v\n", oldPod.Name)
	klog.Infof("new pod: %+v\n", newPod.Name)
}

func onDelete(obj interface{}) {
	pod := obj.(*corev1.Pod)
	klog.Infof("delete a pod: +v", pod.Name)
}

////////////////////////////////////////////////////////////////////////////////
var NodeNameIndex = "nodeName"

// NodeNameIndexFunc 按主机名称对 pod 列表进行索引查询, 即查询指定主机上的 Pod 列表.
//
// cache缓存, indexer索引器等都是针对单个资源而言的, 不同资源的缓存/索引规则无法通用.
//
// 	@param obj: client-go 会在向 cache 中添加目标资源对象时, 调用索引器函数,
//              截取该对象中的某一字段作为索引键.
//
// 如果把 Pod 看成是一张数据表, NodeName 为表中的一个字段的话, 就相当于对 NodeName 这一列创建索引.
// 而其带来的效果就是, 查询操作时如果以 NodeName 为条件的话, 效率将大大提升, 
// 但同时也会影响增/删速度.
//
// 	@return: 函数器函数的返回值格式是固定的, 必须是一个 []string.
func NodeNameIndexFunc(obj interface{}) ([]string, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return []string{}, nil
	}
	return []string{pod.Spec.NodeName}, nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	klog.Infof("初始化 kube client...")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	klog.Infof("初始化 informer factory...")

	// Shared指的是多个 lister 共享同一个cache, 而且资源的变化会同时通知到cache和listers.
	factory := informers.NewSharedInformerFactory(clientset, 60*time.Second)

	// podInformer 拥有两个方法: Informer, Lister.
	// 其实可以把 Informer 看作是 watch 操作.
	podInformer := factory.Core().V1().Pods()
	nodeInformer := factory.Core().V1().Nodes()
	defer runtime.HandleCrash()

	myIndexers := cache.Indexers{
		// namespace 是 client-go 内置的默认索引器.
		// 但是这里不能写, 因为会和原来的造成冲突(AddIndexers()实现中会与其合并)
		// cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
		NodeNameIndex:        NodeNameIndexFunc,
	}
	err = podInformer.Informer().AddIndexers(myIndexers)
	if err != nil {
		klog.Errorf("添加自定义索引器失败: %s", err)
	}

	// 启动 informer, 开始 list & watch 流程(不需要使用 go func() 模式)
	factory.Start(stopCh)
	// 从 apiserver 同步某种资源的全部对象.
	// 之后就可以使用watch这种资源, 维护这份缓存.
	factory.WaitForCacheSync(stopCh)

	// 使用自定义 handler, 处理 watch 响应的各种事件.
	// 具体的维护操作在informer内部执行, 这里挂载的是额外的触发器.
	// 需要注意的是, 在上面的list过程中, 会不断触发onAdd事件, 相当于服务发现了.
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	// 从informer对象创建lister, 不过这里的代码没有特殊的目的,
	// 应该只是展示一下通过informer的接口得到list资源的方法.
	// 从 lister 中获取所有 items
	nodeList, err := nodeInformer.Lister().List(labels.Everything())
	if err != nil {
		klog.Errorf("获取Node列表失败: %s", err)
	}
	klog.Infof("获取Node列表成功")
	for _, node := range nodeList {
		klog.Infof("Node: %s ======================", node.Name)
		podsOnTheNode, err := podInformer.Informer().GetIndexer().ByIndex(NodeNameIndex, node.Name)
		if err != nil {
			klog.Errorf("查询Node: %s 索引失败: %s", node.Name, err)
		}
		for _, pod := range podsOnTheNode {
			klog.Infof(pod.(*corev1.Pod).Name)
		}
	}

	<-stopCh
}
