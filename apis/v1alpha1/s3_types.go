package v1alpha1

type S3Region string

type S3SignatureVersion string

const (
	S3SignatureVersionV2 = "s3v2"
	S3SignatureVersionV4 = "s3v4"
)

type S3Protocol struct {
	Endpoint         string
	BucketName       string
	Region           string
	SignatureVersion S3SignatureVersion
}
