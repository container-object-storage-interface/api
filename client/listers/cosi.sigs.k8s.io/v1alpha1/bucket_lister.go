package v1alpha1

import (
	bucketv1alpha1 "github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type bucketLister struct {
	indexer cache.Indexer
}

func NewBucketLister(indexer cache.Indexer) BucketLister {
	return &bucketLister{
		indexer: indexer,
	}
}

type BucketLister interface {
	List(labels.Selector) ([]*bucketv1alpha1.Bucket, error)
	Get(string) (*bucketv1alpha1.Bucket, error)

	BucketListerExpansion
}

func (b *bucketLister) List(selector labels.Selector) (ret []*bucketv1alpha1.Bucket, err error) {
	err = cache.ListAll(b.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*bucketv1alpha1.Bucket))
	})
	return ret, err
}

func (b *bucketLister) Get(name string) (*bucketv1alpha1.Bucket, error) {
	obj, exists, err := b.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(bucketv1alpha1.Resource("Bucket"), name)
	}
	return obj.(*bucketv1alpha1.Bucket), nil
}

type BucketListerExpansion interface{}
