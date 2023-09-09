# rsp

本项目实现了对目标集群中 deployment 资源的增删改查操作, 尤其是 labels 信息的修改, 提供了 http 接口.

## 环境搭建

```
docker run -d --name remote-java --privileged=true -p 10001:22 -p 10080:80 -v /usr/local/maven-m2:/root/.m2 registry.cn-hangzhou.aliyuncs.com/generals-space/remote-java:8
```

## 编译构建

```
mvn package -DskipTests
```

## 知识点

- spring boot 工程
- kubernetes java client
- vscode 远程开发环境
- mvn package 构建 jar 包
- kube资源(deployment)更新方式, patch接口的使用
- 自定义json响应体 `ResponseData{status, message, data}`
- http接口全局异常捕获
- lombok 注解精简 getter/setter 方法
