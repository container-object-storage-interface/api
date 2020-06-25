package cosiv1alpha1

import (
	"k8s.io/client-go/rest"
)

type Interface interface {
	RESTClient() rest.Interface
	
	Bucket() BucketV1alpha1Interface
	BucketClass() BucketClassV1Alpha1Interface
	BucketContent() BucketContentV1Alpha1Interface
}

type cosiV1alpha1Client struct {
	restClient *rest.Client

	bucket        *bucketV1alpha1Client
	bucketClass   *bucketClassV1alpha1Client
	bucketContent *bucketContentV1alpha1Client
}

func New(c *rest.Client) *cosiV1alpha1Client {
	return &cosiV1alpha1Client{
		restClient:    c,
		bucket:        newBucketClient(c),
		bucketClass:   newBucketClassClient(c),
		bucketContent: newBucketContentClient(c),
	}
}

func (c *cosiV1alpha1Client) Bucket() *bucketV1alpha1Client {
	return c.bucket
}

func (c *cosiV1alpha1Client) BucketClass() *bucketClassV1alpha1Client {
	return c.bucketClass
}

func (c *cosiV1alpha1Client) BucketContent() *bucketContentV1alpha1Client {
	return c.bucketContent
}

func (c *cosiV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
