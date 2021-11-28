# kube-operator

参考文章

1. [client-go系列之4---Indexer](https://zhuanlan.zhihu.com/p/266512431)
    - 其中 2.1.1 节的索引器数据结构值得一看.
2. [client-go Indexer索引器](https://herbguo.gitbook.io/client-go/informer#4.2-indexer-suo-yin-qi)
    - index 的使用示例...有点难懂
3. [client-go 之 Indexer 的理解](https://blog.51cto.com/u_15077560/2584555)
    - 这篇文章的 indexer 示例代码比参考文章2的易懂.

kuber: 1.16.2

本示例实现了按主机名称对 pod 列表进行索引查询, 即查询指定主机上的 Pod 列表.

> cache缓存, indexer索引器等都是针对单个资源而言的, 不同资源的缓存/索引规则无法通用.

client-go 会在向 cache 中添加目标资源对象时, 调用索引器函数, 截取该对象中的某一字段作为索引键.

如果把 Pod 看成是一张数据表, NodeName 为表中的一个字段的话, 就相当于对 NodeName 这一列创建索引.

而其带来的效果就是, 查询操作时如果以 NodeName 为条件的话, 效率将大大提升, 但同时也会影响增/删速度.

## 运行

工程路径放置在`XXX/generals-space/kube-operator`

```console
$ go run main.go
I1128 20:38:46.972581    4236 main.go:63] 初始化 kube client...
I1128 20:38:46.995989    4236 main.go:72] 初始化 informer factory...
I1128 20:38:47.601481    4236 main.go:120] 获取Node列表成功
I1128 20:38:47.607322    4236 main.go:122] Node: ly-xjf-r021702-gyt ======================
I1128 20:38:47.638086    4236 main.go:128] kube-controller-manager-ly-xjf-r021702-gyt
I1128 20:38:47.679592    4236 main.go:128] kube-scheduler-ly-xjf-r021702-gyt
I1128 20:38:47.682520    4236 main.go:128] cni-delivery-2b7q5
I1128 20:38:47.755765    4236 main.go:128] etcd-ly-xjf-r021702-gyt
I1128 20:38:47.798734    4236 main.go:128] kube-apiserver-ly-xjf-r021702-gyt
I1128 20:38:47.812405    4236 main.go:128] coredns-74676bfdf5-49z74
I1128 20:38:47.818752    4236 main.go:128] kube-proxy-llv96
I1128 20:38:47.820706    4236 main.go:122] Node: ly-xjf-r021703-gyt ======================
I1128 20:38:47.822170    4236 main.go:128] cni-delivery-5qnpb
I1128 20:38:47.837796    4236 main.go:128] kube-proxy-n5zcs
I1128 20:38:47.899808    4236 main.go:122] Node: ly-xjf-r021704-gyt ======================
I1128 20:38:47.910550    4236 main.go:128] cni-delivery-9xgh4
I1128 20:38:47.943785    4236 main.go:128] coredns-flink-776bb66868-fcwjd
I1128 20:38:47.945252    4236 main.go:128] kube-proxy-5shwz
```
