package v1alpha1

import (
	"context"
	"time"

	bucketv1alpha1 "github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1"
	bucketv1alpha1lister "github.com/container-object-storage-interface/api/client/listers/cosi.sigs.k8s.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	kubernetes "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"

	"github.com/container-object-storage-interface/api/client/clientset"
)

func NewBucketRequestInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBucketRequestInformer(client, namespace, resyncPeriod, indexers, nil)
}

func NewFilteredBucketRequestInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CosiV1alpha1().BucketRequests(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CosiV1alpha1().BucketRequests(namespace).Watch(context.TODO(), options)
			},
		},
		&bucketv1alpha1.BucketRequest{},
		resyncPeriod,
		indexers,
	)
}

type BucketRequestInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() bucketv1alpha1lister.BucketRequestLister
}

type bucketRequestInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace string
}

func (f *bucketRequestInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&bucketv1alpha1.BucketRequest{}, f.defaultInformer)
}

func (f *bucketRequestInformer) Lister() bucketv1alpha1lister.BucketRequestLister {
	return bucketv1alpha1lister.NewBucketRequestLister(f.Informer().GetIndexer())
}

func (f *bucketRequestInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	cl, ok := client.(clientset.Interface)
	if !ok {
		panic("Invalid client for bucketrequest informer passed. Must be cosiv1alpha1.Interface, found kubernetes.Interface instead")
	}
	return NewFilteredBucketRequestInformer(cl,
		f.namespace,
		resyncPeriod,
		cache.Indexers{
			cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
		}, f.tweakListOptions)
}
