# kube-operator

kuber: 1.16.2

分别使用`informers.NewSharedInformerFactory()`和`cache.NewSharedIndexInformer()`两个方法, 对`Service`和`Pod`资源进行监听, 实现的效果基本相同.

## 运行

工程路径放置在`XXX/generals-space/kube-operator`

```console
$ go run main.go
I0330 18:14:15.093787   36338 main.go:55] 初始化 informer...
I0330 18:14:17.194311   36338 main.go:25] add a node:k8s-master-01
I0330 18:14:17.194299   36338 main.go:93] 获取主机列表: [&Node{ObjectMeta:{k8s-master-01   /api/v1/nodes/k8s-master-01...
```
