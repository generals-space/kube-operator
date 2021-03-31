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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	kubegroupv1 "generals-space/kube-operator/pkg/apis/kubegroup/v1"
	versioned "generals-space/kube-operator/pkg/client/clientset/versioned"
	internalinterfaces "generals-space/kube-operator/pkg/client/informers/externalversions/internalinterfaces"
	v1 "generals-space/kube-operator/pkg/client/listers/kubegroup/v1"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PodClusterInformer provides access to a shared informer and lister for
// PodClusters.
type PodClusterInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.PodClusterLister
}

type podClusterInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPodClusterInformer constructs a new informer for PodCluster type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPodClusterInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPodClusterInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPodClusterInformer constructs a new informer for PodCluster type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPodClusterInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubegroupV1().PodClusters(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubegroupV1().PodClusters(namespace).Watch(options)
			},
		},
		&kubegroupv1.PodCluster{},
		resyncPeriod,
		indexers,
	)
}

func (f *podClusterInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPodClusterInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *podClusterInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kubegroupv1.PodCluster{}, f.defaultInformer)
}

func (f *podClusterInformer) Lister() v1.PodClusterLister {
	return v1.NewPodClusterLister(f.Informer().GetIndexer())
}