package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BucketPhase string

const (
	BucketPhasePending BucketPhase = "pending"
	BucketPhaseBound   BucketPhase = "bound"
)

type BucketRequestBinding struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type BucketRequestSpec struct {
	// +optional
	BucketName      string            `json:"bucketName,omitempty"`
	// +optional
	SecretName      string            `json:"secretName,omitempty"`
	// +optional
	BucketPrefix    string            `json:"bucketPrefix,omitempty"`
	// +optional
	BucketClassName string            `json:"bucketClassName,omitempty"`
	Protocol        ProtocolSignature `json:"protocol"`
}

type BucketRequestStatus struct {
	Phase BucketPhase `json:"phase,omitempty"`
}

type AnonymousAccessMode struct {
	Private         bool `json:"private,omitempty"`
	PublicReadOnly  bool `json:"publicReadOnly,omitempty"`
	PublicReadWrite bool `json:"publicReadWrite,omitempty"`
}

type BucketSpec struct {
	Provisioner         string              `json:"provisioner"`
	AnonymousAccessMode AnonymousAccessMode `json:"anonymousAccessMode,omitempty"`
	BucketClassName     string              `json:"bucketClassName,omitempty"`
	// +listType=set
	PermittedNamespaces []string            `json:"permittedNamespaces,omitempty"`
	Protocol            Protocol            `json:"protocol"`
	Parameters          map[string]string   `json:"parameters,omitempty"`
}

type BucketStatus struct {
	Message             string               `json:"message,omitempty"`
	Phase               BucketPhase          `json:"phase,omitempty"`
	BoundBucketRequests BucketRequestBinding `json:"boundBucketRequests,omitempty"`
}

type ReleasePolicy string

const (
	ReleasePolicyRetain ReleasePolicy = "retain"
	ReleasePolicyDelete ReleasePolicy = "delete"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Namespaced,path=bucketRequests
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

type BucketRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BucketRequestSpec   `json:"spec,omitempty"`
	Status BucketRequestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketRequest `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,path=buckets
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

type Bucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BucketSpec   `json:"spec,omitempty"`
	Status BucketStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bucket `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,path=bucketClasses
// +kubebuilder:storageversion

type BucketClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Provisioner                   string                `json:"provisioner,omitempty"`
	IsDefaultBucketClass          bool                  `json:"isDefaultBucketClass,omitempty"`
	// +listType=set
	AdditionalPermittedNamespaces []string              `json:"additionalPermittedNamespaces,omitempty"`

	// +listType=atomic
	SupportedProtocols            []Protocol            `json:"supportedProtocols"`

	// +listType=set
	AnonymousAccessModes          []AnonymousAccessMode `json:"anonymousAccessModes,omitempty"`
	ReleasePolicy                 ReleasePolicy         `json:"releasePolicy,omitempty"`
	Parameters                    map[string]string     `json:"parameters,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketClass `json:"items"`
}

// Bucket Access Types

type PolicyActions struct {
	// +listType=set
	Allow []string `json:"allow,omitempty"`
	// +listType=set
	Deny  []string `json:"deny,omitempty"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,path=bucketAccessClasses
// +kubebuilder:storageversion

type BucketAccessClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Provisioner        string              `json:"provisioner,omitempty"`
	PolicyActions      PolicyActions       `json:"policyActions,omitempty"`

	// +listType=set
	SupportedProtocols []ProtocolSignature `json:"supportedProtocols,omitempty"`

	Parameters         map[string]string   `json:"parameters,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketAccessClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketAccessClass `json:"items"`
}

type BucketAccessSpec struct {
	BucketAccessRequestName      string            `json:"bucketAccessRequestName,omitempty"`
	BucketAccessRequestNamespace string            `json:"bucketAccessRequestNamespace,omitempty"`
	ServiceAccountName           string            `json:"serviceAccountName,omitempty"`
	KeySecretName                string            `json:"keySecretName,omitempty"`
	Provisioner                  string            `json:"provisioner,omitempty"`
	Parameters                   map[string]string `json:"parameters,omitempty"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,path=bucketAccesses
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

type BucketAccess struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BucketAccessSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketAccessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketAccess `json:"items"`
}

type BucketAccessRequestSpec struct {
	ServiceAccountName    string `json:"serviceAccountName,omitempty"`
	AccessSecretName      string `json:"accessSecretName,omitempty"`
	BucketRequestName     string `json:"bucketRequestName,omitempty"`
	BucketAccessClassName string `json:"bucketAccessClassName,omitempty"`
	BucketAccessName      string `json:"bucketAccessName,omitempty"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Namespaced,path=bucketAccessRequests
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

type BucketAccessRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BucketAccessRequestSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketAccessRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketAccessRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bucket{}, &BucketList{})
	SchemeBuilder.Register(&BucketRequest{}, &BucketRequestList{})
	SchemeBuilder.Register(&BucketClass{}, &BucketClassList{})

	SchemeBuilder.Register(&BucketAccess{}, &BucketAccessList{})
	SchemeBuilder.Register(&BucketAccessRequest{}, &BucketAccessRequestList{})
	SchemeBuilder.Register(&BucketAccessClass{}, &BucketAccessClassList{})
}
