# 启动manager失败-if kind is a CRD, it should be installed before calling Start

```console
$ /podcluster_manager --kubeconfig /etc/kubernetes/admin.conf
{"level":"info","ts":1617890594.9171607,"logger":"entrypoint","msg":"setting up client for manager"}
...省略
{"level":"info","ts":1617890597.0275443,"logger":"controller-runtime.manager","msg":"starting metrics server","path":"/metrics"}
{"level":"info","ts":1617890597.0276124,"logger":"controller-runtime.controller","msg":"Starting EventSource","controller":"podcluster-controller","source":"kind source: /, Kind="}
{"level":"error","ts":1617890600.625955,"logger":"controller-runtime.source","msg":"if kind is a CRD, it should be installed before calling Start","kind":"PodCluster.kubegroup.generals.space","error":"no matches for kind \"PodCluster\" in version \"kubegroup.generals.space/v1\"","stacktrace":"github.com/go-logr/zapr.(*zapLogger).Error\n\t/usr/local/gopath/pkg/mod/github.com/go-logr/zapr@v0.1.1/zapr.go:128\nsigs.k8s.io/controller-runtime/pkg/source.(*Kind).Start\n\t/usr/local/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.6.0/pkg/source/source.go:105\nsigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func1\n\t/usr/local/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.6.0/pkg/internal/controller/controller.go:165\nsigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start\n\t/usr/local/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.6.0/pkg/internal/controller/controller.go:198\nsigs.k8s.io/controller-runtime/pkg/manager.(*controllerManager).startLeaderElectionRunnables.func1\n\t/usr/local/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.6.0/pkg/manager/internal.go:514"}
{"level":"error","ts":1617890600.6262786,"logger":"entrypoint","msg":"unable to run the manager","error":"no matches for kind \"PodCluster\" in version \"kubegroup.generals.space/v1\"","stacktrace":"github.com/go-logr/zapr.(*zapLogger).Error\n\t/usr/local/gopath/pkg/mod/github.com/go-logr/zapr@v0.1.1/zapr.go:128\nmain.main\n\t/home/generals-space/kube-operator/cmd/manager/main.go:66\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}
```

需要事先创建CRD资源.
