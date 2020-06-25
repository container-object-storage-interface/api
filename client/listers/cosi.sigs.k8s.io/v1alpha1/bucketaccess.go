package v1alpha1

import (
	bucketv1alpha1 "github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type bucketAccessLister struct {
	indexer cache.Indexer
}

func NewBucketAccessLister(indexer cache.Indexer) BucketAccessLister {
	return &bucketAccessLister{
		indexer: indexer,
	}
}

type BucketAccessLister interface {
	List(labels.Selector) ([]*bucketv1alpha1.BucketAccess, error)
	Get(string) (*bucketv1alpha1.BucketAccess, error)

	BucketAccessListerExpansion
}

func (b *bucketAccessLister) List(selector labels.Selector) (ret []*bucketv1alpha1.BucketAccess, err error) {
	err = cache.ListAll(b.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*bucketv1alpha1.BucketAccess))
	})
	return ret, err
}

func (b *bucketAccessLister) Get(name string) (*bucketv1alpha1.BucketAccess, error) {
	obj, exists, err := b.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(bucketv1alpha1.Resource("BucketAccess"), name)
	}
	return obj.(*bucketv1alpha1.BucketAccess), nil
}

type BucketAccessListerExpansion interface{}
