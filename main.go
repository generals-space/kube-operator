package main

import (
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	internalapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	kubeletconfiginternal "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/kubelet/dockershim"
	dockerremote "k8s.io/kubernetes/pkg/kubelet/dockershim/remote"
	"k8s.io/kubernetes/pkg/kubelet/remote"
	"k8s.io/kubernetes/pkg/kubelet/server/streaming"
)

// getRuntimeAndImageServices 创建到 dockershim(grpc server)的客户端对象.
//
// 	@param remoteRuntimeEndpoint: unix:///var/run/dockershim.sock
// 	@param remoteImageEndpoint: unix:///var/run/dockershim.sock
func getRuntimeAndImageServices(
	remoteRuntimeEndpoint, remoteImageEndpoint string,
	runtimeRequestTimeout metav1.Duration,
) (rs internalapi.RuntimeService, is internalapi.ImageManagerService, err error) {
	rs, err = remote.NewRemoteRuntimeService(
		remoteRuntimeEndpoint, runtimeRequestTimeout.Duration,
	)
	if err != nil {
		return
	}
	is, err = remote.NewRemoteImageService(
		remoteImageEndpoint, runtimeRequestTimeout.Duration,
	)
	if err != nil {
		return
	}
	return
}

func main() {
	// pkg/kubelet/kubelet.go -> NewMainKubelet()

	pluginSettings := dockershim.NetworkPluginSettings{
		HairpinMode: kubeletconfiginternal.HairpinMode( // --hairpin-mode
			"promiscuous-bridge",
		),
		NonMasqueradeCIDR:  "10.0.0.0/8",         // --non-masquerade-cidr
		PluginName:         "cni",                // --network-plugin
		PluginConfDir:      "/etc/cni/net.d",     // --cni-conf-dir
		PluginBinDirString: "/opt/cni/bin",       // --cni-bin-dir
		PluginCacheDir:     "/var/lib/cni/cache", // --cni-cache-dir
		MTU:                1460,                 // --network-plugin-mtu
	}

	// 用于调用 docker 官方库, 创建与 dockerd 服务进行通信的客户端对象.
	dockerClientConfig := &dockershim.ClientConfig{
		// --docker-endpoint
		DockerEndpoint: "unix:///var/run/docker.sock",
		// --runtime-request-timeout="2m0s"
		RuntimeRequestTimeout: time.Second * 120,
		// --image-pull-progress-deadline="1m0s"
		ImagePullProgressDeadline: time.Second * 60,
	}
	// 这个也是用于与 dockerd 服务通信的, 主要实现 exec, logs 等请求.
	//
	// Create and start the CRI shim running as a grpc server.
	streamingConfig := &streaming.Config{
		// --streaming-connection-idle-timeout="4h0m0s"
		StreamIdleTimeout:               time.Hour * 4,
		StreamCreationTimeout:           streaming.DefaultConfig.StreamCreationTimeout,
		SupportedRemoteCommandProtocols: streaming.DefaultConfig.SupportedRemoteCommandProtocols,
		SupportedPortForwardProtocols:   streaming.DefaultConfig.SupportedPortForwardProtocols,
	}

	// kubernetes-v1.17.2 cmd/kubelet/app/options/container_runtime.go
	podSandboxImage := "registry.cn-hangzhou.aliyuncs.com/google_containers/pause/pause:3.1"
	// --runtime-cgroups=""
	runtimeCgroups := ""
	// --cgroup-driver="systemd"
	cgroupDriver := "systemd"
	// --experimental-dockershim-root-directory="/var/lib/dockershim"
	dockershimRootDirectory := "/var/lib/dockershim"
	// --redirect-container-streaming="false"
	redirectContainerStreaming := false
	ds, err := dockershim.NewDockerService(
		dockerClientConfig,
		podSandboxImage,
		streamingConfig,
		&pluginSettings,
		runtimeCgroups,
		cgroupDriver,
		dockershimRootDirectory,
		!redirectContainerStreaming,
	)

	if err != nil {
		log.Printf("failed to create docker service: %s", err)
		return
	}

	// dockershim.sock 并不是 docker 自身生成的文件,
	// 而是 kubelet 在启动时创建 grpc server 时创建的.
	// --container-runtime-endpoint="unix:///var/run/dockershim.sock"
	dockershimEP := "unix:///var/run/dockershim.sock"

	// docker server 启动时, 会建立与 dockerd 服务(/var/run/docker.sock)的连接对象,
	// 同时启动一个名为 dockershim grpc 服务, 该服务只提供2个 Service, 
	// 即为下面的 runtimeService 与 imageService, 分别用于容器/镜像的操作.
	//
	// kubelet 调用这两个 Service 时, dockershim 会将请求转换为到 dockerd 的格式.
	server := dockerremote.NewDockerServer(dockershimEP, ds)
	if err := server.Start(); err != nil {
		log.Printf("failed to start docker server: %s", err)
		return
	}

	// --runtime-request-timeout="2m0s"
	runtimeRequestTimeout := metav1.Duration{time.Second * 120}
	runtimeService, imageService, err := getRuntimeAndImageServices(
		dockershimEP, dockershimEP, runtimeRequestTimeout,
	)
	if err != nil {
		log.Printf("failed to init grpc service: %s", err)
		return
	}
	log.Println("============================= containers")

	// 查看所有容器
	containers, err := runtimeService.ListContainers(&runtimeapi.ContainerFilter{})
	if err != nil {
		log.Printf("failed to list containers: %s", err)
		return
	}
	for _, c := range containers {
		log.Printf("container: %s", c.Metadata.Name)
	}
	log.Println("============================= images")
	// 查看所有镜像
	images, err := imageService.ListImages(&runtimeapi.ImageFilter{})
	if err != nil {
		log.Printf("failed to list images: %s", err)
		return
	}
	for _, i := range images {
		log.Printf("container: %+v", i.RepoTags)
	}
}
