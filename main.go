package main

import (
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"generals-space/kube-operator/pkg/controller/pod"
	"generals-space/kube-operator/pkg/controller/svc"
)

func main() {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	// 初始化 client
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}
	stopCh := make(chan struct{})
	defer close(stopCh)

	// 用 factory 的方式得到 informer(NewSharedInformerFactory())
	// factory 中其实可以包含 n 个 indexInformer, 不过ta们的 cache 是共享的.
	svcController := svc.New(kubeClient)
	svcController.Start(stopCh)

	// 用 index 的方式得到 informer(NewSharedIndexInformer())
	podController := pod.New(kubeClient)
	podController.Start(stopCh)

	<-stopCh
}
