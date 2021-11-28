package main

import (
	"fmt"
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

/*
	各种处理函数中, obj即为watch接口响应得到的资源对象.
*/

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

func main() {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	// 初始化 client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	klog.Infof("初始化 informer...")
	// Shared指的是多个 lister 共享同一个cache, 而且资源的变化会同时通知到cache和listers.
	factory := informers.NewSharedInformerFactory(clientset, 60 * time.Second)

	// podInformer 拥有两个方法: Informer, Lister.
	// 其实可以把 Informer 看作是 watch 操作.
	podInformer := factory.Core().V1().Pods()
	informer := podInformer.Informer()
	defer runtime.HandleCrash()

	// 启动 informer, 开始 list & watch 流程(不需要使用 go func() 模式)
	factory.Start(stopCh)

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
	podLister := podInformer.Lister()
	// 从 lister 中获取所有 items
	podList, err := podLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("获取Pod列表失败: %s", err)
	}
	klog.Infof("获取Pod列表")
	for _, pod := range podList {
		klog.Infof("%s", pod.Name)
	}

	<-stopCh
}
