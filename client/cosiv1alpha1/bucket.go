package cosiv1alpha1

import (
	cosiapiv1alpha1 "github.com/container-object-storage-interface/api/apis/v1alpha1"
)

type BucketV1alpha1Interface interface {
	Create(*cosiapiv1alpha1.Bucket) (*cosiapiv1alpha1.Bucket, error)
	Update(*cosiapiv1alpha1.Bucket) (*cosiapiv1alpha1.Bucket, error)
	UpdateStatus(*cosiapiv1alpha1.Bucket) (*cosiapiv1alpha1.Bucket, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*cosiapiv1alpha1.Bucket, error)
	List(opts metav1.ListOptions) (*cosiapiv1alpha1.BucketList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *cosiapiv1alpha1.Bucket, err error)
}

