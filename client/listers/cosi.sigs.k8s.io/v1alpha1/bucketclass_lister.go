package v1alpha1

import (
	bucketv1alpha1 "github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type bucketClassLister struct {
	indexer cache.Indexer
}

func NewBucketClassLister(indexer cache.Indexer) BucketClassLister {
	return &bucketClassLister{
		indexer: indexer,
	}
}

type BucketClassLister interface {
	List(labels.Selector) ([]*bucketv1alpha1.BucketClass, error)
	Get(string) (*bucketv1alpha1.BucketClass, error)

	BucketClassListerExpansion
}

func (b *bucketClassLister) List(selector labels.Selector) (ret []*bucketv1alpha1.BucketClass, err error) {
	err = cache.ListAll(b.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*bucketv1alpha1.BucketClass))
	})
	return ret, err
}

func (b *bucketClassLister) Get(name string) (*bucketv1alpha1.BucketClass, error) {
	obj, exists, err := b.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(bucketv1alpha1.Resource("BucketClass"), name)
	}
	return obj.(*bucketv1alpha1.BucketClass), nil
}

type BucketClassListerExpansion interface{}
