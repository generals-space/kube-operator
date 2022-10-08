package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"k8s.io/klog"
	appsv1 "k8s.io/api/apps/v1"
	apimYaml "k8s.io/apimachinery/pkg/util/yaml"
	cgKube "k8s.io/client-go/kubernetes"
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
		klog.Errorf("failed to build kubeconfig: %s", err.Error())
		return
	}

	kubeClient, err := cgKube.NewForConfig(cfg)
	if err != nil {
		klog.Errorf("failed to get kube client: %s", err.Error())
		return
	}

	sts := &appsv1.StatefulSet{}
	yamlbytes, err := ioutil.ReadFile("sts.yaml")
	if err != nil {
		klog.Errorf("failed to read sts file: %s", err.Error())
		return
	}

	reader := bytes.NewReader(yamlbytes)
	decoder := apimYaml.NewYAMLOrJSONDecoder(reader, len(yamlbytes))
	err = decoder.Decode(sts)
	if err != nil {
		klog.Errorf("failed to parse sts file: %s", err.Error())
		return
	}
	klog.Infof("%+v\n", sts)

	stsObj, err := kubeClient.AppsV1().StatefulSets("default").Create(sts)
	if err != nil {
		klog.Errorf("failed to create sts: %s", err.Error())
		return
	}
	klog.Infof("=== sts: %+v\n", stsObj)
}
