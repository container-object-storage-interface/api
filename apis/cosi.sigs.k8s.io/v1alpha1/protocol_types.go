package v1alpha1

type ProtocolSignature string

const (
	ProtocolSignatureS3    ProtocolSignature = "s3"
	ProtocolSignatureAzure ProtocolSignature = "azure"
	ProtcolSignatureGCS    ProtocolSignature = "gcs"
)

type Protocol struct {
	ProtocolSignature ProtocolSignature `json:"protocolSignature"`
	S3                S3Protocol        `json:"s3,omitempty"`
	Azure             AzureProtocol     `json:"azure,omitempty"`
	GCS               GCSProtocol       `json:"gcs,omitempty"`
}
