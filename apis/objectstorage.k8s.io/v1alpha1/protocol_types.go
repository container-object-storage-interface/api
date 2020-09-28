package v1alpha1

type ProtocolName string

const (
	ProtocolNameS3    ProtocolName = "s3"
	ProtocolNameAzure ProtocolName = "azure"
	ProtocolNameGCS   ProtocolName = "gcs"
)

type Protocol struct {
	// +kubebuilder:validation:Enum:={s3,azure,gcs}
	ProtocolName ProtocolName `json:"protocolName"`
	// +optional
	S3 *S3Protocol `json:"s3,omitempty"`
	// +optional
	AzureBlob *AzureProtocol `json:"azureBlob,omitempty"`
	// +optional
	GCS *GCSProtocol `json:"gcs,omitempty"`
}
