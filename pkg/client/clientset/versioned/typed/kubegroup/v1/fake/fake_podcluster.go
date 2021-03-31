/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	kubegroupv1 "generals-space/kube-operator/pkg/apis/kubegroup/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePodClusters implements PodClusterInterface
type FakePodClusters struct {
	Fake *FakeKubegroupV1
	ns   string
}

var podclustersResource = schema.GroupVersionResource{Group: "kubegroup.generals.space", Version: "v1", Resource: "podclusters"}

var podclustersKind = schema.GroupVersionKind{Group: "kubegroup.generals.space", Version: "v1", Kind: "PodCluster"}

// Get takes name of the podCluster, and returns the corresponding podCluster object, and an error if there is any.
func (c *FakePodClusters) Get(name string, options v1.GetOptions) (result *kubegroupv1.PodCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(podclustersResource, c.ns, name), &kubegroupv1.PodCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubegroupv1.PodCluster), err
}

// List takes label and field selectors, and returns the list of PodClusters that match those selectors.
func (c *FakePodClusters) List(opts v1.ListOptions) (result *kubegroupv1.PodClusterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(podclustersResource, podclustersKind, c.ns, opts), &kubegroupv1.PodClusterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &kubegroupv1.PodClusterList{ListMeta: obj.(*kubegroupv1.PodClusterList).ListMeta}
	for _, item := range obj.(*kubegroupv1.PodClusterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested podClusters.
func (c *FakePodClusters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(podclustersResource, c.ns, opts))

}

// Create takes the representation of a podCluster and creates it.  Returns the server's representation of the podCluster, and an error, if there is any.
func (c *FakePodClusters) Create(podCluster *kubegroupv1.PodCluster) (result *kubegroupv1.PodCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(podclustersResource, c.ns, podCluster), &kubegroupv1.PodCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubegroupv1.PodCluster), err
}

// Update takes the representation of a podCluster and updates it. Returns the server's representation of the podCluster, and an error, if there is any.
func (c *FakePodClusters) Update(podCluster *kubegroupv1.PodCluster) (result *kubegroupv1.PodCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(podclustersResource, c.ns, podCluster), &kubegroupv1.PodCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubegroupv1.PodCluster), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePodClusters) UpdateStatus(podCluster *kubegroupv1.PodCluster) (*kubegroupv1.PodCluster, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(podclustersResource, "status", c.ns, podCluster), &kubegroupv1.PodCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubegroupv1.PodCluster), err
}

// Delete takes name of the podCluster and deletes it. Returns an error if one occurs.
func (c *FakePodClusters) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(podclustersResource, c.ns, name), &kubegroupv1.PodCluster{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePodClusters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(podclustersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &kubegroupv1.PodClusterList{})
	return err
}

// Patch applies the patch and returns the patched podCluster.
func (c *FakePodClusters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *kubegroupv1.PodCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(podclustersResource, c.ns, name, pt, data, subresources...), &kubegroupv1.PodCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubegroupv1.PodCluster), err
}
