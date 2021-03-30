package main

import (
	"path/filepath"
	"time"

	"k8s.io/klog"
	cgKuber "k8s.io/client-go/kubernetes"
	cgClientcmd "k8s.io/client-go/tools/clientcmd"
	cgHomedir "k8s.io/client-go/util/homedir"

	clientset "generals-space/kube-operator/pkg/client/clientset/versioned"
	crdInformerFactory "generals-space/kube-operator/pkg/client/informers/externalversions"
	"generals-space/kube-operator/pkg/signals"
)

func main() {
	// 处理信号
	stopCh := signals.SetupSignalHandler()

	// 先尝试从 ~/.kube 目录下获取配置, 如果没有, 则尝试寻找 Pod 内置的认证配置
	var kubeconfig string
	home := cgHomedir.HomeDir()
	kubeconfig = filepath.Join(home, ".kube", "config")
	cfg, err := cgClientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Errorf("Error building kubeconfig: %s", err.Error())
	}

	// kubeClient 用于集群内资源操作, crdClient 用于操作 crd 资源本身.
	// 具体区别目前还不清楚, 不过示例中大多都是这么做的.
	kubeClient, err := cgKuber.NewForConfig(cfg)
	if err != nil {
		klog.Errorf("Error building kubernetes clientset: %s", err.Error())
	}
	crdClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Errorf("Error building example clientset: %s", err.Error())
	}

	crdInformerFactory := crdInformerFactory.NewSharedInformerFactory(crdClient, time.Second*30)

	//得到controller
	controller := NewController(
		kubeClient,
		crdClient,
		crdInformerFactory.Kubegroup().V1().PodClusters(),
	)

	//启动informer
	go crdInformerFactory.Start(stopCh)

	//controller开始处理消息
	if err = controller.Run(2, stopCh); err != nil {
		klog.Errorf("Error running controller: %s", err.Error())
	}
}
