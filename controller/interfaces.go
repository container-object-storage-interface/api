package controller

import (
	"context"

	// storage
	"github.com/container-object-storage-interface/api/apis/storage.k8s.io/v1alpha1"
	bucketclientset "github.com/container-object-storage-interface/api/clientset"
	
	// k8s client
	kubeclientset "k8s.io/client-go/kubernetes"
)

// Set the clients for each of the listeners
type GenericListener interface {
	InitializeKubeClient(kubeclientset.Interface)
	InitializeBucketClient(bucketclientset.Interface)
}

type BucketListener interface {
	GenericListener

	Add(ctx context.Context, b *v1alpha1.Bucket) error
	Update(ctx context.Context, old *v1alpha1.Bucket, new *v1alpha1.Bucket) error
	Delete(ctx context.Context, b *v1alpha1.Bucket) error
}

func (c *Controller) AddBucketListener(b BucketListener) {
	c.initialized = true
	c.BucketListener = b
}

type BucketClassListener interface {
	GenericListener

	Add(ctx context.Context, b *v1alpha1.BucketClass) error
	Update(ctx context.Context, old *v1alpha1.BucketClass, new *v1alpha1.BucketClass) error
	Delete(ctx context.Context, b *v1alpha1.BucketClass) error
}

func (c *Controller) AddBucketClassListener(b BucketClassListener) {
	c.initialized = true
	c.BucketClassListener = b
}

type BucketRequestListener interface {
	GenericListener

	Add(ctx context.Context, b *v1alpha1.BucketRequest) error
	Update(ctx context.Context, old *v1alpha1.BucketRequest, new *v1alpha1.BucketRequest) error
	Delete(ctx context.Context, b *v1alpha1.BucketRequest) error
}

func (c *Controller) AddBucketRequestListener(b BucketRequestListener) {
	c.initialized = true
	c.BucketRequestListener = b
}

type BucketAccessListener interface {
	GenericListener

	Add(ctx context.Context, b *v1alpha1.BucketAccess) error
	Update(ctx context.Context, old *v1alpha1.BucketAccess, new *v1alpha1.BucketAccess) error
	Delete(ctx context.Context, b *v1alpha1.BucketAccess) error
}

func (c *Controller) AddBucketAccessListener(b BucketAccessListener) {
	c.initialized = true
	c.BucketAccessListener = b
}

type BucketAccessClassListener interface {
	GenericListener

	Add(ctx context.Context, b *v1alpha1.BucketAccessClass) error
	Update(ctx context.Context, old *v1alpha1.BucketAccessClass, new *v1alpha1.BucketAccessClass) error
	Delete(ctx context.Context, b *v1alpha1.BucketAccessClass) error
}

func (c *Controller) AddBucketAccessClassListener(b BucketAccessClassListener) {
	c.initialized = true
	c.BucketAccessClassListener = b
}

type BucketAccessRequestListener interface {
	GenericListener

	Add(ctx context.Context, b *v1alpha1.BucketAccessRequest) error
	Update(ctx context.Context, old *v1alpha1.BucketAccessRequest, new *v1alpha1.BucketAccessRequest) error
	Delete(ctx context.Context, b *v1alpha1.BucketAccessRequest) error
}

func (c *Controller) AddBucketAccessRequestListener(b BucketAccessRequestListener) {
	c.initialized = true
	c.BucketAccessRequestListener = b
}
