module generals-space/kube-operator

go 1.12

require (
	github.com/elazarl/goproxy v0.0.0-20180725130230-947c36da3153 // indirect
	github.com/googleapis/gnostic v0.1.0 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

replace (
	k8s.io/apimachinery => /home/k8s.io/apimachinery
	k8s.io/client-go => /home/k8s.io/client-go
)
