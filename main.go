package main

import (
	"path/filepath"

	"k8s.io/klog"
	apieClientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apimMetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cgKuber "k8s.io/client-go/kubernetes"
	cgClientcmd "k8s.io/client-go/tools/clientcmd"
	cgHomedir "k8s.io/client-go/util/homedir"
)

func main() {
	// 先尝试从 ~/.kube 目录下获取配置, 如果没有, 则尝试寻找 Pod 内置的认证配置
	var kubeconfig string
	home := cgHomedir.HomeDir()
	kubeconfig = filepath.Join(home, ".kube", "config")
	cfg, err := cgClientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	// kubeClient 用于集群内资源操作, crdClient 用于操作 crd 资源本身.
	// 具体区别目前还不清楚, 不过示例中大多都是这么做的.
	kubeClient, err := cgKuber.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}
	deploy, err := kubeClient.AppsV1().Deployments("kube-system").Get("coredns", apimMetav1.GetOptions{})
	if err != nil {
		klog.Errorf("get coredns deploy failed: %s", err)
		return
	}
	klog.Infof("%+v\n", deploy)

	apieClient, err := apieClientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building apiextensions clientset: %s", err.Error())
		return
	}
	crds, err := apieClient.ApiextensionsV1().CustomResourceDefinitions().List(apimMetav1.ListOptions{})
	if err != nil {
		klog.Errorf("list crd failed: %s", err)
		return
	}
	for _, crd := range crds.Items {
		klog.Infof("%s\n", crd.Name)
	}
}
