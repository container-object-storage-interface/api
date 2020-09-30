package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type BucketRequestBinding struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type BucketRequestSpec struct {
	// +optional
	BucketInstanceName string `json:"bucketInstanceName,omitempty"`
	// +optional
	BucketPrefix string `json:"bucketPrefix,omitempty"`
	// +optional
	BucketClassName string    `json:"bucketClassName,omitempty"`
	Protocol        Protocol  `json:"protocol"`
}

type BucketRequestStatus struct {
	// +optional
	Message string `json:"message,omitempty"`
	// +optional
	BucketAvailable bool `json:"bucketAvailable,omitempty"`
}

type AnonymousAccessMode struct {
	// +optional
	Private bool `json:"private,omitempty"`
	// +optional
	PublicReadOnly bool `json:"publicReadOnly,omitempty"`
	// +optional
	PublicReadWrite bool `json:"publicReadWrite,omitempty"`
	// +optional
	PublicWriteOnly bool `json:"publicWriteOnly,omitempty"`
}

type BucketSpec struct {
	Provisioner string `json:"provisioner"`
	// +kubebuilder:default:=retain
	RetentionPolicy RetentionPolicy `json:"retentionPolicy"`

	AnonymousAccessMode AnonymousAccessMode `json:"anonymousAccessMode,omitempty"`
	BucketClassName     string              `json:"bucketClassName,omitempty"`
	BucketRequest       *ObjectReference    `json:"bucketRequest,omitempty"`
	// +listType=atomic
	AllowedNamespaces []string `json:"allowedNamespaces,omitempty"`

	Protocol Protocol `json:"protocol"`
	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

type BucketStatus struct {
	// +optional
	Message string `json:"message,omitempty"`
	// +optional
	BucketAvailable bool `json:"bucketAvailable,omitempty"`
}

type RetentionPolicy string

const (
	RetentionPolicyRetain RetentionPolicy = "Retain"
	RetentionPolicyDelete RetentionPolicy = "Delete"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

type BucketRequest struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BucketRequestSpec `json:"spec,omitempty"`
	// +optional
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
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

type Bucket struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BucketSpec `json:"spec,omitempty"`
	// +optional
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
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

type BucketClass struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Provisioner string `json:"provisioner,omitempty"`
	// +optional
	// +kubebuilder:default:=false
	IsDefaultBucketClass bool `json:"isDefaultBucketClass,omitempty"`
	// +listType=atomic
	// +optional
	AllowedNamespaces []string `json:"allowedNamespaces,omitempty"`
	// +optional
	Protocol string `json:"protocol"`
	// +listType=atomic
	AnonymousAccessModes []AnonymousAccessMode `json:"anonymousAccessModes,omitempty"`
	// +kubebuilder:default:=retain
	RetentionPolicy RetentionPolicy `json:"retentionPolicy,omitempty"`
	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketClass `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

type BucketAccessClass struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Provisioner string `json:"provisioner,omitempty"`

	PolicyActions string `json:"policyActions,omitempty"`

	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketAccessClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketAccessClass `json:"items"`
}

type BucketAccessSpec struct {
	Provisioner string `json:"provisioner,omitempty"`

	// +optional
	BucketInstanceName string `json:"bucketInstanceName,omitempty"`
	// +optional
	BucketAccessRequest string `json:"bucketAccessRequest,omitempty"`
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// +optional
	AccessSecretName string `json:"accessSecretName,omitempty"`
	// +optional
	MintedSecretName string `json:"mintedSecretName,omitempty"`
	// +optional
	PolicyActionsConfigMapData map[string]string `json:"policyActionsConfigMapData,omitempty"`
	// +optional
	Principal string `json:"principal,omitempty"`
	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

type BucketAccessStatus struct {
	// +optional
	Message string `json:"message,omitempty"`
	// +optional
	AccessGranted bool `json:"accessGranted,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

type BucketAccess struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BucketAccessSpec `json:"spec,omitempty"`
	// +optional
	Status BucketAccessStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketAccessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketAccess `json:"items"`
}

type BucketAccessRequestSpec struct {
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// +optional
	AccessSecretName  string `json:"accessSecretName,omitempty"`
	BucketRequestName string `json:"bucketRequestName,omitempty"`
	// +optional
	BucketAccessClassName string `json:"bucketAccessClassName,omitempty"`
	// +optional
	BucketAccessName string `json:"bucketAccessName,omitempty"`
}

type BucketAccessRequestStatus struct {
	// +optional
	Message string `json:"message,omitempty"`
	// +optional
	AccessGranted bool `json:"accessGranted,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

type BucketAccessRequest struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BucketAccessRequestSpec `json:"spec,omitempty"`
	// +optional
	Status BucketAccessRequestStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BucketAccessRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BucketAccessRequest `json:"items"`
}

type ObjectReference struct {
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty"`
	// UID of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids
	// +optional
	UID types.UID `json:"uid,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Bucket{}, &BucketList{})
	SchemeBuilder.Register(&BucketRequest{}, &BucketRequestList{})
	SchemeBuilder.Register(&BucketClass{}, &BucketClassList{})

	SchemeBuilder.Register(&BucketAccess{}, &BucketAccessList{})
	SchemeBuilder.Register(&BucketAccessRequest{}, &BucketAccessRequestList{})
	SchemeBuilder.Register(&BucketAccessClass{}, &BucketAccessClassList{})
}
