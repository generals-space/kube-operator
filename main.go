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

func getRuntimeAndImageServices(
	remoteRuntimeEndpoint string,
	remoteImageEndpoint string,
	runtimeRequestTimeout metav1.Duration,
) (internalapi.RuntimeService, internalapi.ImageManagerService, error) {
	rs, err := remote.NewRemoteRuntimeService(
		remoteRuntimeEndpoint,
		runtimeRequestTimeout.Duration,
	)
	if err != nil {
		return nil, nil, err
	}
	is, err := remote.NewRemoteImageService(
		remoteImageEndpoint,
		runtimeRequestTimeout.Duration,
	)
	if err != nil {
		return nil, nil, err
	}
	return rs, is, err
}

func main() {
	// pkg/kubelet/kubelet.go -> NewMainKubelet()

	pluginSettings := dockershim.NetworkPluginSettings{
		// --hairpin-mode
		HairpinMode: kubeletconfiginternal.HairpinMode("promiscuous-bridge"),
		// --non-masquerade-cidr
		NonMasqueradeCIDR: "10.0.0.0/8",
		// --network-plugin
		PluginName: "cni",
		// --cni-conf-dir
		PluginConfDir: "/etc/cni/net.d",
		// --cni-bin-dir
		PluginBinDirString: "/opt/cni/bin",
		// --cni-cache-dir
		PluginCacheDir: "/var/lib/cni/cache",
		// --network-plugin-mtu
		MTU: 1460,
	}

	// Create and start the CRI shim running as a grpc server.
	streamingConfig := &streaming.Config{
		// --streaming-connection-idle-timeout="4h0m0s"
		StreamIdleTimeout:               time.Hour * 4,
		StreamCreationTimeout:           streaming.DefaultConfig.StreamCreationTimeout,
		SupportedRemoteCommandProtocols: streaming.DefaultConfig.SupportedRemoteCommandProtocols,
		SupportedPortForwardProtocols:   streaming.DefaultConfig.SupportedPortForwardProtocols,
	}
	dockerClientConfig := &dockershim.ClientConfig{
		// --docker-endpoint
		DockerEndpoint: "unix:///var/run/docker.sock",
		// --runtime-request-timeout="2m0s"
		RuntimeRequestTimeout: time.Second * 120,
		// --image-pull-progress-deadline="1m0s"
		ImagePullProgressDeadline: time.Second * 60,
	}
	// kubernetes-v1.17.2 cmd/kubelet/app/options/container_runtime.go
	podSandboxImage := "k8s.gcr.io/pause:3.1"
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

	server := dockerremote.NewDockerServer(dockershimEP, ds)
	if err := server.Start(); err != nil {
		log.Printf("failed to start docker server: %s", err)
		return
	}

	// --runtime-request-timeout="2m0s"
	runtimeRequestTimeout := metav1.Duration{time.Second * 120}
	runtimeService, imageService, err := getRuntimeAndImageServices(
		dockershimEP,
		dockershimEP,
		runtimeRequestTimeout,
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
