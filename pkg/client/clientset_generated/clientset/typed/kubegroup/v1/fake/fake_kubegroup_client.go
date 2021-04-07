// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "generals-space/kube-operator/pkg/client/clientset_generated/clientset/typed/kubegroup/v1"

	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeKubegroupV1 struct {
	*testing.Fake
}

func (c *FakeKubegroupV1) PodClusters(namespace string) v1.PodClusterInterface {
	return &FakePodClusters{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKubegroupV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
