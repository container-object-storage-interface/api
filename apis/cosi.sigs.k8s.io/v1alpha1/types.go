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
	Bucket          string            `json:"bucket"`
	SecretName      string            `json:"secretName,omitempty"`
	Provisioner     string            `json:"provisioner,omitempty"`
	BucketPrefix    string            `json:"bucketPrefix,omitempty"`
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

type BucketClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Provisioner                   string                `json:"provisioner,omitempty"`
	IsDefaultBucketClass          bool                  `json:"isDefaultBucketClass,omitempty"`
	AdditionalPermittedNamespaces []string              `json:"additionalPermittedNamespaces,omitempty"`
	SupportedProtocols            []Protocol            `json:"supportedProtocols"`
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

func init() {
	SchemeBuilder.Register(&Bucket{}, &BucketList{})
	SchemeBuilder.Register(&BucketRequest{}, &BucketRequestList{})
	SchemeBuilder.Register(&BucketClass{}, &BucketClassList{})
}
