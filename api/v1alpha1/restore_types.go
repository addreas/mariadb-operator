package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RestoreSpec defines the desired state of restore
type RestoreSpec struct {
	// ContainerTemplate defines templates to configure Container objects.
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ContainerTemplate `json:",inline"`
	// PodTemplate defines templates to configure Pod objects.
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PodTemplate `json:",inline"`

	// RestoreSource defines a source for restoring a MariaDB.
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RestoreSource `json:",inline"`
	// MariaDBRef is a reference to a MariaDB object.
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MariaDBRef MariaDBRef `json:"mariaDbRef" webhook:"inmutable"`
	// LogLevel to be used n the Backup Job. It defaults to 'info'.
	// +optional
	// +kubebuilder:default=info
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	LogLevel string `json:"logLevel,omitempty"`
	// BackoffLimit defines the maximum number of attempts to successfully perform a Backup.
	// +optional
	// +kubebuilder:default=5
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	BackoffLimit int32 `json:"backoffLimit,omitempty"`
	// RestartPolicy to be added to the Backup Job.
	// +optional
	// +kubebuilder:default=OnFailure
	// +kubebuilder:validation:Enum=Always;OnFailure;Never
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RestartPolicy corev1.RestartPolicy `json:"restartPolicy,omitempty" webhook:"inmutable"`
}

// RestoreStatus defines the observed state of restore
type RestoreStatus struct {
	// Conditions for the Restore object.
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors={"urn:alm:descriptor:io.kubernetes.conditions"}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (r *RestoreStatus) SetCondition(condition metav1.Condition) {
	if r.Conditions == nil {
		r.Conditions = make([]metav1.Condition, 0)
	}
	meta.SetStatusCondition(&r.Conditions, condition)
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=rmdb
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Complete",type="string",JSONPath=".status.conditions[?(@.type==\"Complete\")].status"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Complete\")].message"
// +kubebuilder:printcolumn:name="MariaDB",type="string",JSONPath=".spec.mariaDbRef.name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +operator-sdk:csv:customresourcedefinitions:resources={{Restore,v1alpha1},{Job,v1}}

// Restore is the Schema for the restores API. It is used to define restore jobs and its restoration source.
type Restore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RestoreSpec   `json:"spec,omitempty"`
	Status RestoreStatus `json:"status,omitempty"`
}

func (r *Restore) IsComplete() bool {
	return meta.IsStatusConditionTrue(r.Status.Conditions, ConditionTypeComplete)
}

// +kubebuilder:object:root=true

// RestoreList contains a list of restore
type RestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Restore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Restore{}, &RestoreList{})
}
