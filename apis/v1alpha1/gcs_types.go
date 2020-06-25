package v1alpha1

type GCSProtocol struct {
	BucketName     string `json:"bucketName"`
	PrivateKeyName string `json:"privateKeyName"`
	ProjectID      string `json:"projectID"`
	ServiceAccount string `json:"serviceAccount"`
}
