/*
Copyright 2019 caicloud authors. All rights reserved.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	kubernetes "github.com/caicloud/clientset/kubernetes"
	v1alpha1 "github.com/caicloud/clientset/listers/workload/v1alpha1"
	workloadv1alpha1 "github.com/caicloud/clientset/pkg/apis/workload/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	clientgokubernetes "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
)

// WorkloadRevisionInformer provides access to a shared informer and lister for
// WorkloadRevisions.
type WorkloadRevisionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.WorkloadRevisionLister
}

type workloadRevisionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewWorkloadRevisionInformer constructs a new informer for WorkloadRevision type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewWorkloadRevisionInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredWorkloadRevisionInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredWorkloadRevisionInformer constructs a new informer for WorkloadRevision type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredWorkloadRevisionInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.WorkloadV1alpha1().WorkloadRevisions(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.WorkloadV1alpha1().WorkloadRevisions(namespace).Watch(options)
			},
		},
		&workloadv1alpha1.WorkloadRevision{},
		resyncPeriod,
		indexers,
	)
}

func (f *workloadRevisionInformer) defaultInformer(client clientgokubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredWorkloadRevisionInformer(client.(kubernetes.Interface), f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *workloadRevisionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&workloadv1alpha1.WorkloadRevision{}, f.defaultInformer)
}

func (f *workloadRevisionInformer) Lister() v1alpha1.WorkloadRevisionLister {
	return v1alpha1.NewWorkloadRevisionLister(f.Informer().GetIndexer())
}
