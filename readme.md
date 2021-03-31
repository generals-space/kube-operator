# kube-operator

## 前言

生成的代码已经加入代码库, 如要验证生成步骤, 可见上一个commit.

------

本项目使用[code-generator](https://github.com/kubernetes/code-generator)生成CRD的代码(不先生成代码的话本项目无法直接运行), 需要事先准备好当前的目录结构, 主要为

1. `pkg/apis/kubegroup/v1/doc.go`
2. `pkg/apis/kubegroup/v1/register.go`
3. `pkg/apis/kubegroup/v1/types.go`

这3个文件的内容可以看作是确定的格式, `doc.go`与`registry.go`可以直接忽略, 在定义自己的CRD时, 需要修改`doc.go`与`registry.go`中的`group`名称, 以及CRD结构体对象, 除了具体的类型定义, 其他的像`import`的内容, `+genclient`和`+k8s`这种编译标记, 也都拷贝过去.

我们主要关心`types.go`文件, `types.go`声明了我们CRD对象的结构(本例中为`PodCluster`), 而在这个结构体中, operator 工程一般只关注`Spec`与`Status`成员, 为了简便, `Spec`下只定义了一个成员属性: `PodReplicas`.

本示例中, `group`定为`kubegroup.generals.space`, `version`定为`v1`.

> `pkg/apis/kubegroup/v1/doc.go`中的`kubegroup`即为`group`名称中的第1个段.

在使用`code-generator`生成代码前, 只要这3个文件就可以了, 当前工程的`main.go`, `controller.go`以及`pkg/signals`都是生成代码后再编写的.

## 生成代码(code-generator)

- code-generator: tag v0.17.2
- apimachinery: tag v0.17.2

### 环境准备

首先, 需要`code-generator`和`apimachinery`两个工程存在于`$GOPATH/src/k8s.io/`目录下, `go mod`形式的依赖管理无效(放在CRD工程的`vendor`目录中根本没用). 否则在执行脚本时会出现`Hit an unsupported type invalid type for invalid type`的问题.

> 使用`GO111MODULE=off go get -v`或是直接使用`git clone`都行, 不过需要注意将两个工程的分支切换到`tag v0.17.2`, 然后使用`go mod vendor`安装ta们两个各自的依赖.

```bash
mkdir -p $GOPATH/src/k8s.io
cd $GOPATH/src/k8s.io/

git clone https://github.com/kubernetes/code-generator.git
cd code-generator
git checkout -b v0.17.2 v0.17.2
go mod vendor
cd ..

git clone https://github.com/kubernetes/apimachinery.git
cd apimachinery
git checkout -b v0.17.2 v0.17.2
go mod vendor
cd ..
```

> 注意: 上面我们使用`git`将两个项目clone到本地, 因为`apimachinery`必须要在`GOPATH`目录下, 否则在生成过程中可能出现`Hit an unsupported type invalid type for invalid type`.

### 代码生成

理论上, 使用`code-generator`生成代码时, 要求我们的工程也放在`GOPATH`下, 否则生成会失败.

不过现在很多工程都是用`go mod`构建了, 我们需要把工程拷贝到`GOPATH`下, 生成代码后再拷出来(当你想再生成一种CRD资源时, 还得再搞一个来回)...

有没有比较好一点的方法呢? 

其实可以建立当前工程路径到`GOPATH`下的软链接

```
$ mkdir -p $GOPATH/src/generals-space
$ pwd
/home/generals-space/kube-operator
$ ln -s /home/generals-space/kube-operator $GOPATH/src/generals-space/kube-operator
```

然后执行如下命令开始生成代码

```
$GOPATH/src/k8s.io/code-generator/generate-groups.sh all generals-space/kube-operator/pkg/client generals-space/kube-operator/pkg/apis kubegroup:v1
```

执行命令时所在的目录没有强制要求, `kubegroup:v1`应该指定了`pkg/apis`下的CRD源.

执行完成后, `pkg`目录下会新增`client`目录, 用于操作我们的CRD对象.

## 启动

kuber: 1.17.2

需要先创建crd资源, `kubectl apply -f deploy/01.crd.yaml`

```
$ go run *.go
I0331 00:04:07.012887   63088 controller.go:59] Creating event broadcaster
I0331 00:04:07.013372   63088 controller.go:80] Setting up event handlers
I0331 00:04:07.014235   63088 controller.go:104] 开始controller业务, 开始一次缓存数据同步
I0331 00:04:07.115051   63088 controller.go:109] worker启动
I0331 00:04:07.115096   63088 controller.go:114] worker已经启动
^CI0331 00:04:15.082303   63088 controller.go:116] worker已经结束
```

------

如果出现如下报错, 则说明未事先创建`CRD`资源

```
$ go run *.go
I0331 00:03:01.289358   62397 controller.go:59] Creating event broadcaster
I0331 00:03:01.289687   62397 controller.go:80] Setting up event handlers
I0331 00:03:01.289716   62397 controller.go:104] 开始controller业务, 开始一次缓存数据同步
E0331 00:03:03.307087   62397 reflector.go:153] pkg/mod/k8s.io/client-go@v0.17.2/tools/cache/reflector.go:105: Failed to list *v1.PodCluster: the server could not find the requested resource (get podclusters.kubegroup.generals.space)
```
