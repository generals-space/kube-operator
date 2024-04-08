# kube-operator

1. [01.simple-informer](../../tree/01.simple-informer)
    - 最简 controller, 单文件, 不需要创建 CRD 资源和对应的 golang 对象(包含 GVK, Spec, Status等), 只监听 kuber 内置的资源对象(本工程中为Node主机资源)的变动;
2. [02.crd-podcluster](../../tree/02.crd-podcluster)
    - 声明`PodCluster`类型的CRD资源, 通过`code-generator`生成代码, 附详细的操作方法;
3. [03.crd-podcluster-kubebuilder](../../tree/03.crd-podcluster-kubebuilder)
    - 声明`PodCluster`类型的CRD资源, 通过`kubebuilder`生成代码, 附详细的操作方法;
    - `Reconcile()`主方法;
    - `kustomize`工具生成`yaml`文件;
4. [04.informer-factory](../../tree/04.informer-factory)
    - `informers.NewSharedInformerFactory()`和`cache.NewSharedIndexInformer()`两个方法, 对Service和Pod资源进行监听, 实现的效果基本相同;
    - 简要描述了两种`informer`的关系;
5. [05.node-pod-indexer](../../tree/05.node-pod-indexer)
    - 介绍了Indexer索引器的原理及使用方法;
    - 实现了按主机名称对 pod 列表进行索引查询, 即查询指定主机上的 Pod 列表;
6. [rest-client](../../tree/rest-client)
    - client-go rest-client 的使用方法, 直接构造 http 请求, 自主选择目标资源路径;
    - 与 clientset 有所区别;
7. [create-from-yaml](../../tree/create-from-yaml)
    - 读入 yaml 文件, 构造 Object 对象创建资源;
8. [http-watch-api](../../tree/http-watch-api)
    - 实现了trunk接口的 server 端, 可以通过 curl 该接口实现 watch 的效果;
    - 普通的 http server, 不依赖任何外部库;
9. [apiextensions-apiserver-client](../../tree/apiextensions-apiserver-client)
    - 实现了`kubectl get crd`的功能(与常规内置资源, 及开发者自定义的crd资源有所不同)
10. [docker-api](../../tree/docker-api)
    - 实现了`docker ps`, `docker images`的功能;
    - 基于 kubelet 引用到的[kubernetes-sigs/cri-tools](https://github.com/kubernetes-sigs/cri-tools)库实现;
    - dockershim grpc 服务调用方式;
11. [java-client-patch-deployment](../../tree/java-client-patch-deployment)
    - 实现了对目标集群中 deployment 资源的增删改查操作, 尤其是 labels 信息的修改, 提供了 http 接口.
    - spring boot 工程
    - kubernetes java client
    - vscode 远程开发环境
    - mvn package 构建 jar 包
    - kube资源(deployment)更新方式, patch接口的使用
    - 自定义json响应体 `ResponseData{status, message, data}`
    - http接口全局异常捕获
    - lombok 注解精简 getter/setter 方法
12. [node-pod-indexer-kubebuilder](../../tree/node-pod-indexer-kubebuilder)
    - 实现了在 kubebuilder 工程中, 按主机名称对 pod 列表进行索引查询, 即查询指定主机上的 Pod 列表;
    - 功能上等同于示例 05.
