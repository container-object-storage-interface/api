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
	cache "k8s.io/client-go/tools/cache"
	kubernetes "k8s.io/client-go/kubernetes"
	
	"github.com/container-object-storage-interface/api/client/clientset"
)

func NewBucketClassInformer(client clientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBucketClassInformer(client, resyncPeriod, indexers, nil)
}

func NewFilteredBucketClassInformer(client clientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CosiV1alpha1().BucketClasses().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CosiV1alpha1().BucketClasses().Watch(context.TODO(), options)
			},
		},
		&bucketv1alpha1.BucketClass{},
		resyncPeriod,
		indexers,
	)
}

type BucketClassInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() bucketv1alpha1lister.BucketClassLister
}

type bucketClassInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

func (f *bucketClassInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&bucketv1alpha1.BucketClass{}, f.defaultInformer)
}

func (f *bucketClassInformer) Lister() bucketv1alpha1lister.BucketClassLister {
	return bucketv1alpha1lister.NewBucketClassLister(f.Informer().GetIndexer())
}

func (f *bucketClassInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	cl, ok := client.(clientset.Interface)
	if !ok {
		panic("Invalid client for bucketclass informer passed. Must be cosiv1alpha1.Interface, found kubernetes.Interface instead")
	}
	return NewFilteredBucketClassInformer(cl, resyncPeriod, cache.Indexers{
		cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
	}, f.tweakListOptions)
}
