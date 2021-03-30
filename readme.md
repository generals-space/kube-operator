# kube-operator

kuber: 1.16.2

可以说, 这个示例就是最简的 controller 了, 不需要创建 CRD 资源和对应的 golang 对象(包含 GVK, Spec, Status等), 只监听 kuber 内置的资源对象(本工程中为`Node`主机资源)的变动.

`informer`存在的意义就是, 在kuber集群各组件在与apiserver进行通信时添加一个中间缓存层, 减轻apiserver的负载压力. `informer`内置缓存功能, 可保证与apiserver查询到的数据保持一致.

## 运行

工程路径放置在`XXX/generals-space/kube-operator`

```console
$ go run main.go
I0330 18:14:15.093787   36338 main.go:55] 初始化 informer...
I0330 18:14:17.194311   36338 main.go:25] add a node:k8s-master-01
I0330 18:14:17.194299   36338 main.go:93] 获取主机列表: [&Node{ObjectMeta:{k8s-master-01   /api/v1/nodes/k8s-master-01...
```
