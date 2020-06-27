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

func NewBucketAccessRequestInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBucketAccessRequestInformer(client, namespace, resyncPeriod, indexers, nil)
}

func NewFilteredBucketAccessRequestInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CosiV1alpha1().BucketAccessRequests(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CosiV1alpha1().BucketAccessRequests(namespace).Watch(context.TODO(), options)
			},
		},
		&bucketv1alpha1.BucketAccessRequest{},
		resyncPeriod,
		indexers,
	)
}

type BucketAccessRequestInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() bucketv1alpha1lister.BucketAccessRequestLister
}

type bucketAccessRequestInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace string
}

func (f *bucketAccessRequestInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&bucketv1alpha1.BucketAccessRequest{}, f.defaultInformer)
}

func (f *bucketAccessRequestInformer) Lister() bucketv1alpha1lister.BucketAccessRequestLister {
	return bucketv1alpha1lister.NewBucketAccessRequestLister(f.Informer().GetIndexer())
}

func (f *bucketAccessRequestInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	cl, ok := client.(clientset.Interface)
	if !ok {
		panic("Invalid client for bucketaccessRequest informer passed. Must be cosiv1alpha1.Interface, found kubernetes.Interface instead")
	}
	return NewFilteredBucketAccessRequestInformer(cl,
		f.namespace,
		resyncPeriod,
		cache.Indexers{
			cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
		}, f.tweakListOptions)
}
