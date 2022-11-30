# kube-operator

参考文章

1. [kubernetes/apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver)

实现类似`kubectl get crds`的功能.

kube 内置的资源对象, 使用的是通过`client-go`构建的客户端`clientset`/`restclient`, 但是这两个客户端无法操作CRD类型. 

开发者如果希望在代码中动态创建/查找/删除 CRD 类型的资源, 需要借助`apiextensions-apiserver`工具库来实现.

当然, 这只是CRD资源, 开发者自定义的资源不叫CRD, 应该叫CR. 操作CR则需要使用通过`code-generator`生成的clientset库了.
