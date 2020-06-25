package v1alpha1

import (
	bucketv1alpha1 "github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type bucketAccessRequestLister struct {
	indexer cache.Indexer
}

func NewBucketAccessRequestLister(indexer cache.Indexer) BucketAccessRequestLister {
	return &bucketAccessRequestLister{
		indexer: indexer,
	}
}

type BucketAccessRequestLister interface {
	List(labels.Selector) ([]*bucketv1alpha1.BucketAccessRequest, error)
	BucketAccessRequests(string) BucketAccessRequestNamespaceLister

	BucketAccessRequestListerExpansion
}

func (b *bucketAccessRequestLister) List(selector labels.Selector) (ret []*bucketv1alpha1.BucketAccessRequest, err error) {
	err = cache.ListAll(b.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*bucketv1alpha1.BucketAccessRequest))
	})
	return ret, err
}

func (b *bucketAccessRequestLister) BucketAccessRequests(namespace string) BucketAccessRequestNamespaceLister {
	return &bucketAccessRequestNamespaceLister{indexer: b.indexer, namespace: namespace}
}

type BucketAccessRequestNamespaceLister interface {
	List(labels.Selector) ([]*bucketv1alpha1.BucketAccessRequest, error)
	Get(string) (*bucketv1alpha1.BucketAccessRequest, error)

	BucketAccessRequestNamespaceListerExpansion
}

type bucketAccessRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

func (b *bucketAccessRequestNamespaceLister) List(selector labels.Selector) (ret []*bucketv1alpha1.BucketAccessRequest, err error) {
	err = cache.ListAllByNamespace(b.indexer, b.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*bucketv1alpha1.BucketAccessRequest))
	})
	return ret, err
}

func (b *bucketAccessRequestNamespaceLister) Get(name string) (*bucketv1alpha1.BucketAccessRequest, error) {
	obj, exists, err := b.indexer.GetByKey(b.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(bucketv1alpha1.Resource("BucketAccessRequest"), name)
	}
	return obj.(*bucketv1alpha1.BucketAccessRequest), nil
}

type BucketAccessRequestListerExpansion interface{}
type BucketAccessRequestNamespaceListerExpansion interface{}
