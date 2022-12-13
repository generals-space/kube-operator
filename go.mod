module generals-space/kube-operator

go 1.13

require (
	k8s.io/apimachinery v0.17.2
	k8s.io/cri-api v0.17.2
	k8s.io/kubernetes v1.17.2
)

replace (
	github.com/armon/circbuf => github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e
	github.com/docker/go-connections => github.com/docker/go-connections v0.3.0
	github.com/google/cadvisor => github.com/google/cadvisor v0.35.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.7.0
	github.com/lithammer/dedent => github.com/lithammer/dedent v1.1.0
	github.com/morikuni/aec => github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega => github.com/onsi/gomega v1.7.0
	github.com/opencontainers/go-digest => github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.0-rc9
	github.com/vishvananda/netlink => github.com/vishvananda/netlink v1.0.0
	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery => ../../k8s.io/apimachinery
	k8s.io/apiserver => k8s.io/apiserver v0.17.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.2
	k8s.io/client-go => k8s.io/client-go v0.17.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.2
	k8s.io/code-generator => k8s.io/code-generator v0.17.2
	k8s.io/component-base => k8s.io/component-base v0.17.2
	k8s.io/cri-api => ../../k8s.io/cri-api
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.2
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.2
	k8s.io/kubectl => k8s.io/kubectl v0.17.2
	k8s.io/kubelet => k8s.io/kubelet v0.17.2
	k8s.io/kubernetes => ../../k8s.io/kubernetes
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.2
	k8s.io/metrics => k8s.io/metrics v0.17.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.2
	k8s.io/utils => k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
)
