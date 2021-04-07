


// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

//go:generate deepcopy-gen -O zz_generated.deepcopy -i . -h ../../../../boilerplate.go.txt
//go:generate defaulter-gen -O zz_generated.defaults -i . -h ../../../../boilerplate.go.txt
//go:generate conversion-gen -O zz_generated.conversion -i . -h ../../../../boilerplate.go.txt

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=generals-space/kube-operator/pkg/apis/kubegroup
// +k8s:defaulter-gen=TypeMeta
// +groupName=kubegroup.generals.space
package v1 // import "generals-space/kube-operator/pkg/apis/kubegroup/v1"

