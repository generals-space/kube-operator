# kube-operator

1. [01.simple-informer](../../tree/01.simple-informer)
    - 同名分支
    - 最简 controller, 单文件, 不需要创建 CRD 资源和对应的 golang 对象(包含 GVK, Spec, Status等), 只监听 kuber 内置的资源对象(本工程中为Node主机资源)的变动.
2. [02.crd-podcluster](../../tree/02.crd-podcluster)
    - 同名分支
    - 声明`PodCluster`类型的CRD资源, 通过`code-generator`生成代码, 附详细的操作方法.
3. [03.crd-podcluster-kubebuilder](../../tree/03.crd-podcluster-kubebuilder)
    - 同名分支
    - 声明`PodCluster`类型的CRD资源, 通过`kubebuilder`生成代码, 附详细的操作方法.
    - `Reconcile()`主方法
    - `kustomize`工具生成`yaml`文件
4. [04.informer-factory](../../tree/04.informer-factory)
    - 同名分支
    - `informers.NewSharedInformerFactory()`和`cache.NewSharedIndexInformer()`两个方法, 对Service和Pod资源进行监听, 实现的效果基本相同.
    - 简要描述了两种`informer`的关系.
5. [05.node-pod-indexer](../../tree/05.node-pod-indexer)
    - 同名分支
    - 介绍了Indexer索引器的原理及使用方法.
    - 实现了按主机名称对 pod 列表进行索引查询, 即查询指定主机上的 Pod 列表.
