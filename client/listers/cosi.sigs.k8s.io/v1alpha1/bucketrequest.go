package v1alpha1

import (
	bucketv1alpha1 "github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type bucketRequestLister struct {
	indexer cache.Indexer
}

func NewBucketRequestLister(indexer cache.Indexer) BucketRequestLister {
	return &bucketRequestLister{
		indexer: indexer,
	}
}

type BucketRequestLister interface {
	List(labels.Selector) ([]*bucketv1alpha1.BucketRequest, error)
	BucketRequests(string) BucketRequestNamespaceLister

	BucketRequestListerExpansion
}

func (b *bucketRequestLister) List(selector labels.Selector) (ret []*bucketv1alpha1.BucketRequest, err error) {
	err = cache.ListAll(b.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*bucketv1alpha1.BucketRequest))
	})
	return ret, err
}

func (b *bucketRequestLister) BucketRequests(namespace string) BucketRequestNamespaceLister {
	return &bucketRequestNamespaceLister{indexer: b.indexer, namespace: namespace}
}

type BucketRequestNamespaceLister interface {
	List(labels.Selector) ([]*bucketv1alpha1.BucketRequest, error)
	Get(string) (*bucketv1alpha1.BucketRequest, error)

	BucketRequestNamespaceListerExpansion
}

type bucketRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

func (b *bucketRequestNamespaceLister) List(selector labels.Selector) (ret []*bucketv1alpha1.BucketRequest, err error) {
	err = cache.ListAllByNamespace(b.indexer, b.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*bucketv1alpha1.BucketRequest))
	})
	return ret, err
}

func (b *bucketRequestNamespaceLister) Get(name string) (*bucketv1alpha1.BucketRequest, error) {
	obj, exists, err := b.indexer.GetByKey(b.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(bucketv1alpha1.Resource("BucketRequest"), name)
	}
	return obj.(*bucketv1alpha1.BucketRequest), nil
}

type BucketRequestListerExpansion interface{}
type BucketRequestNamespaceListerExpansion interface{}
