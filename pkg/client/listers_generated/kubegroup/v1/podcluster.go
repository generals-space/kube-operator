// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "generals-space/kube-operator/pkg/apis/kubegroup/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PodClusterLister helps list PodClusters.
type PodClusterLister interface {
	// List lists all PodClusters in the indexer.
	List(selector labels.Selector) (ret []*v1.PodCluster, err error)
	// PodClusters returns an object that can list and get PodClusters.
	PodClusters(namespace string) PodClusterNamespaceLister
	PodClusterListerExpansion
}

// podClusterLister implements the PodClusterLister interface.
type podClusterLister struct {
	indexer cache.Indexer
}

// NewPodClusterLister returns a new PodClusterLister.
func NewPodClusterLister(indexer cache.Indexer) PodClusterLister {
	return &podClusterLister{indexer: indexer}
}

// List lists all PodClusters in the indexer.
func (s *podClusterLister) List(selector labels.Selector) (ret []*v1.PodCluster, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PodCluster))
	})
	return ret, err
}

// PodClusters returns an object that can list and get PodClusters.
func (s *podClusterLister) PodClusters(namespace string) PodClusterNamespaceLister {
	return podClusterNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PodClusterNamespaceLister helps list and get PodClusters.
type PodClusterNamespaceLister interface {
	// List lists all PodClusters in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.PodCluster, err error)
	// Get retrieves the PodCluster from the indexer for a given namespace and name.
	Get(name string) (*v1.PodCluster, error)
	PodClusterNamespaceListerExpansion
}

// podClusterNamespaceLister implements the PodClusterNamespaceLister
// interface.
type podClusterNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all PodClusters in the indexer for a given namespace.
func (s podClusterNamespaceLister) List(selector labels.Selector) (ret []*v1.PodCluster, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PodCluster))
	})
	return ret, err
}

// Get retrieves the PodCluster from the indexer for a given namespace and name.
func (s podClusterNamespaceLister) Get(name string) (*v1.PodCluster, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("podcluster"), name)
	}
	return obj.(*v1.PodCluster), nil
}