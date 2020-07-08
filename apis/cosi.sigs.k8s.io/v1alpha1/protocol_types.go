package v1alpha1

type ProtocolSignature string

const (
	ProtocolSignatureS3    ProtocolSignature = "s3"
	ProtocolSignatureAzure ProtocolSignature = "azure"
	ProtocolSignatureGCS   ProtocolSignature = "gcs"
)

type Protocol struct {
	// +kubebuilder:validation:Enum:={s3,azure,gcs}
	ProtocolSignature ProtocolSignature `json:"protocolSignature"`
	// +optional
	S3 S3Protocol `json:"s3,omitempty"`
	// +optional
	Azure AzureProtocol `json:"azure,omitempty"`
	// +optional
	GCS GCSProtocol `json:"gcs,omitempty"`
}
