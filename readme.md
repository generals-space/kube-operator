# kube-operator

[kubernetes-sigs/cri-tools](https://github.com/kubernetes-sigs/cri-tools)

- golang: 1.13
- kubernetes: v1.17.2

参考kubernetes官方 kubelet 部分源码, 编写的 docker api 接口示例, 分别给出了`docker ps`查看容器列表与`docker images`查看镜像列表的功能.

执行时不需要事先部署 kubelet, 只要有docker即可, 直接运行.

```log
$ go run main.go
I1002 19:36:07.343990    2837 client.go:75] Connecting to docker on unix:///var/run/docker.sock
I1002 19:36:07.344064    2837 client.go:104] Start docker client with request timeout=2m0s
W1002 19:36:07.355998    2837 docker_service.go:563] Hairpin mode set to "promiscuous-bridge" but kubenet is not enabled, falling back to "hairpin-veth"
I1002 19:36:07.356035    2837 docker_service.go:240] Hairpin mode set to "hairpin-veth"
I1002 19:36:07.382055    2837 docker_service.go:255] Docker cri networking managed by cni
I1002 19:36:07.394651    2837 docker_service.go:260] Docker Info: &{ID:RZRX:BCUD:DPGF:M35G:RIU2:TSN7:2T2B:PC5R:D3BC:BO6C:C2C2:CYBT Containers:34 ContainersRunning:1 ContainersPaused:0 ContainersStopped:33 Images:99 Driver:overlay2 DriverStatus:[[Backing Filesystem xfs] [Supports d_type true] [Native Overlay Diff true]] SystemStatus:[] Plugins:{Volume:[local] Network:[bridge host ipvlan macvlan null overlay] Authorization:[] Log:[awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog]} MemoryLimit:true SwapLimit:true KernelMemory:true KernelMemoryTCP:true CPUCfsPeriod:true CPUCfsQuota:true CPUShares:true CPUSet:true PidsLimit:true IPv4Forwarding:true BridgeNfIptables:true BridgeNfIP6tables:true Debug:false NFd:27 OomKillDisable:true NGoroutines:39 SystemTime:2021-10-02T19:36:07.382809996+08:00 LoggingDriver:json-file CgroupDriver:systemd NEventsListener:0 KernelVersion:3.10.0-1160.24.1.el7.x86_64 OperatingSystem:CentOS Linux 7 (Core) OSType:linux Architecture:x86_64 IndexServerAddress:https://index.docker.io/v1/ RegistryConfig:0xc000511c70 NCPU:4 MemTotal:8181817344 GenericResources:[] DockerRootDir:/var/lib/docker HTTPProxy: HTTPSProxy: NoProxy: Name:k8s-master-01 Labels:[] ExperimentalBuild:false ServerVersion:19.03.5 ClusterStore: ClusterAdvertise: Runtimes:map[runc:{Path:runc Args:[]}] DefaultRuntime:runc Swarm:{NodeID: NodeAddr: LocalNodeState:inactive ControlAvailable:false Error: RemoteManagers:[] Nodes:0 Managers:0 Cluster:<nil> Warnings:[]} LiveRestoreEnabled:false Isolation: InitBinary:docker-init ContainerdCommit:{ID:b34a5c8af56e510852c35414db4c1f4fa6172339 Expected:b34a5c8af56e510852c35414db4c1f4fa6172339} RuncCommit:{ID:22c72eb3976b73573b28fe9d14e3f2e113871345-dirty Expected:22c72eb3976b73573b28fe9d14e3f2e113871345-dirty} InitCommit:{ID:fec3683 Expected:fec3683} SecurityOptions:[name=seccomp,profile=default] ProductLicense: Warnings:[]}
I1002 19:36:07.394723    2837 docker_service.go:273] Setting cgroupDriver to systemd
I1002 19:36:07.395024    2837 remote_runtime.go:59] parsed scheme: ""
I1002 19:36:07.395035    2837 remote_runtime.go:59] scheme "" not registered, fallback to default scheme
I1002 19:36:07.395057    2837 passthrough.go:48] ccResolverWrapper: sending update to cc: {[{/var/run/dockershim.sock 0  <nil>}] <nil>}
I1002 19:36:07.395070    2837 clientconn.go:577] ClientConn switching balancer to "pick_first"
I1002 19:36:07.395095    2837 remote_image.go:50] parsed scheme: ""
I1002 19:36:07.395101    2837 remote_image.go:50] scheme "" not registered, fallback to default scheme
I1002 19:36:07.395109    2837 passthrough.go:48] ccResolverWrapper: sending update to cc: {[{/var/run/dockershim.sock 0  <nil>}] <nil>}
I1002 19:36:07.395114    2837 clientconn.go:577] ClientConn switching balancer to "pick_first"
2021-10-02 19:36:07.395124 I | ============================= containers
2021-10-02 19:36:07.398287 I | container: kube-proxy
2021-10-02 19:36:07.398301 I | container: etcd
2021-10-02 19:36:07.398305 I | container: kube-scheduler
2021-10-02 19:36:07.398309 I | container: kube-controller-manager
2021-10-02 19:36:07.398312 I | container: kube-apiserver
2021-10-02 19:36:07.398316 I | container: kube-scheduler
2021-10-02 19:36:07.398319 I | container: kube-apiserver
2021-10-02 19:36:07.398322 I | container: kube-controller-manager
2021-10-02 19:36:07.398325 I | container: etcd
2021-10-02 19:36:07.398328 I | ============================= images
2021-10-02 19:36:07.420746 I | container: [registry.cn-hangzhou.aliyuncs.com/generals-space/centos7-devops:latest]
2021-10-02 19:36:07.420750 I | container: [registry.cn-hangzhou.aliyuncs.com/generals-kuber/crd-ipkeeper:0.0.84]
2021-10-02 19:36:07.420753 I | container: [registry.cn-hangzhou.aliyuncs.com/generals-kuber/cni-terway:0.0.24]
2021-10-02 19:36:07.420778 I | container: [registry.cn-hangzhou.aliyuncs.com/google_containers/kube-proxy:v1.17.2]
2021-10-02 19:36:07.420781 I | container: [registry.cn-hangzhou.aliyuncs.com/google_containers/kube-apiserver:v1.17.2]
2021-10-02 19:36:07.420785 I | container: [registry.cn-hangzhou.aliyuncs.com/google_containers/kube-controller-manager:v1.17.2]
2021-10-02 19:36:07.420788 I | container: [registry.cn-hangzhou.aliyuncs.com/google_containers/kube-scheduler:v1.17.2]
2021-10-02 19:36:07.420792 I | container: [registry.cn-hangzhou.aliyuncs.com/google_containers/coredns:1.6.5]
2021-10-02 19:36:07.420796 I | container: [registry.cn-hangzhou.aliyuncs.com/google_containers/etcd:3.4.3-0]
2021-10-02 19:36:07.420810 I | container: [registry.cn-hangzhou.aliyuncs.com/google_containers/etcd:3.3.15-0]
2021-10-02 19:36:07.420824 I | container: [registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.1]
```
