package main

import (
	"fmt"
	"path/filepath"

	"k8s.io/klog"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	cgScheme "k8s.io/client-go/kubernetes/scheme"
	cgRest "k8s.io/client-go/rest"
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
		return
	}

	// 就是 http 请求头中的 user-agent, 在访问 apiserver 用来标记客户端的访问类型.
	defaultUA := cgRest.DefaultKubernetesUserAgent()
	fmt.Println("default rest user agent: ", defaultUA)
	// cfg.UserAgent 默认为空
	cgRest.AddUserAgent(cfg, "myclient")
	fmt.Println("now rest user agent: ", cfg.UserAgent)

	////////////////////////////////////////////////////////////////////////////////////////////////////
	// 这里的 GroupVersion 决定了访问前缀为 /api/v1
	/*
		// 这种方式也可以
		cfg.ContentConfig = cgRest.ContentConfig{
			GroupVersion:         &corev1.SchemeGroupVersion,
			NegotiatedSerializer: cgScheme.Codecs,
		}
	*/
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	cfg.NegotiatedSerializer = cgScheme.Codecs.WithoutConversion()

	// restClientAPI 请求 /api 接口的 client
	restClientAPI, err := cgRest.RESTClientFor(cfg)
	if err != nil {
		klog.Errorf("failed to create rest client: %s\n", err)
		return
	}
	// 注意这里初始化对象不能是指针, 虽然之后的 Into() 接收的参数是指针, 但ta不会为我们创建目标对象, 
	// 也就不会为这个指针分配空间, 所以我们要事先分配好, 否则会出现空指针报错.
	podList := corev1.PodList{}
	// 实际的请求路径 https://127.0.0.1:6443/api/v1/namespaces/kube-system/pods
	err = restClientAPI.Get().AbsPath("/api/v1").Resource("pods").Namespace("kube-system").Do().Into(&podList)
	if err != nil {
		klog.Errorf("failed to make /api request: %s", err)
		return
	}
	klog.Infof("%+v\n", podList)
	////////////////////////////////////////////////////////////////////////////////////////////////////

	// 这里的 GroupVersion 决定了访问前缀为 /apis/apps/v1
	// 由于 deployment 是在 apps/v1 下的扩展资源, 所以需要这样修改, 尤其是`APIPath` 不要忘记.
	cfg.APIPath = "/apis"
	cfg.ContentConfig = cgRest.ContentConfig{
		GroupVersion:         &appsv1.SchemeGroupVersion,
		NegotiatedSerializer: cgScheme.Codecs,
	}

	// restClientAPIs 请求 /api 接口的 client
	restClientAPIs, err := cgRest.RESTClientFor(cfg)
	if err != nil {
		klog.Errorf("failed to create rest client: %s", err)
		return
	}
	// 注意这里初始化对象不能是指针, 虽然之后的 Into() 接收的参数是指针, 但ta不会为我们创建目标对象, 
	// 也就不会为这个指针分配空间, 所以我们要事先分配好, 否则会出现空指针报错.
	deploys := appsv1.Deployment{}

	getReq := restClientAPIs.Get()
	// getReq.AbsPath("/apis/apps/v1")
	getReq.Resource("deployments")
	getReq.Namespace("kube-system")
	getReq.Name("coredns")
	// 实际的请求路径 https://k8s-server-lb:8443/apis/apps/v1/namespaces/kube-system/deployments/coredns
	klog.Infof("URL: %s\n", getReq.URL().String())
	err = getReq.Do().Into(&deploys)
	if err != nil {
		klog.Errorf("failed to make /apis request: %s", err)
		return
	}
	klog.Infof("%+v\n", deploys)
}
