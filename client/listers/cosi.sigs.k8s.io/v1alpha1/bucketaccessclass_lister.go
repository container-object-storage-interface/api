package v1alpha1

import (
	bucketv1alpha1 "github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type bucketAccessClassLister struct {
	indexer cache.Indexer
}

func NewBucketAccessClassLister(indexer cache.Indexer) BucketAccessClassLister {
	return &bucketAccessClassLister{
		indexer: indexer,
	}
}

type BucketAccessClassLister interface {
	List(labels.Selector) ([]*bucketv1alpha1.BucketAccessClass, error)
	Get(string) (*bucketv1alpha1.BucketAccessClass, error)

	BucketAccessClassListerExpansion
}

func (b *bucketAccessClassLister) List(selector labels.Selector) (ret []*bucketv1alpha1.BucketAccessClass, err error) {
	err = cache.ListAll(b.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*bucketv1alpha1.BucketAccessClass))
	})
	return ret, err
}

func (b *bucketAccessClassLister) Get(name string) (*bucketv1alpha1.BucketAccessClass, error) {
	obj, exists, err := b.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(bucketv1alpha1.Resource("BucketAccessClass"), name)
	}
	return obj.(*bucketv1alpha1.BucketAccessClass), nil
}

type BucketAccessClassListerExpansion interface{}
