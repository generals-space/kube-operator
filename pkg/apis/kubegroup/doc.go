



//go:generate deepcopy-gen -O zz_generated.deepcopy -i . -h ../../../boilerplate.go.txt
//go:generate defaulter-gen -O zz_generated.defaults -i . -h ../../../boilerplate.go.txt

// +k8s:deepcopy-gen=package,register
// +groupName=kubegroup.generals.space

// Package api is the internal version of the API.
package kubegroup

