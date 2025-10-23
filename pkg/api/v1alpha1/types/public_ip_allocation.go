package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AllocationHistory represents a historical allocation entry.
type AllocationHistory struct {
	ContainerInterface string      `json:"containerInterface"`
	Address            string      `json:"address"`     // IP address used
	MacAddress         string      `json:"macAddress"`  // MAC used
	Owner              string      `json:"owner"`       // owner kind/namespace/name
	AllocatedAt        metav1.Time `json:"allocatedAt"` // when allocated
	ReleasedAt         metav1.Time `json:"releasedAt"`  // when released
}

// PublicIPAllocationSpec defines the desired allocation information.
type PublicIPAllocationSpec struct {
	PublicIpPoolName   string `json:"publicIpPoolName"` // Name of PublicIPPool this IP came from
	ContainerInterface string `json:"containerInterface"`
	IpFamily           string `json:"ipFamily"`
	Address            string `json:"address"`    // Allocated IP
	MacAddress         string `json:"macAddress"` // Allocated MAC
	ResourceKind       string `json:"resourceKind"`
	ResourceNamespace  string `json:"resourceNamespace"`
	ResourceName       string `json:"resourceName"`
}

// PublicIPAllocationStatus defines the observed allocation state and history.
type PublicIPAllocationStatus struct {
	AllocatedAt metav1.Time         `json:"allocatedAt"`          // When it was allocated
	ReleasedAt  *metav1.Time        `json:"releasedAt,omitempty"` // When it was released
	History     []AllocationHistory `json:"history,omitempty"`    // Past allocations for audit/history
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=pipa
// +kubebuilder:printcolumn:name="Pool",type=string,JSONPath=`.spec.poolRef`
// +kubebuilder:printcolumn:name="Address",type=string,JSONPath=`.spec.address`
// +kubebuilder:printcolumn:name="MAC",type=string,JSONPath=`.spec.macAddress`
// +kubebuilder:printcolumn:name="Owner",type=string,JSONPath=`.spec.ownerRef.name`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.status.active`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// PublicIPAllocation is the Schema for the publicipallocations API.
type PublicIPAllocation struct {
	metav1.TypeMeta    `json:",inline"`
	*metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PublicIPAllocationSpec   `json:"spec"`
	Status PublicIPAllocationStatus `json:"status"`
}

// +kubebuilder:object:root=true

// PublicIPAllocationList contains a list of PublicIPAllocation.
type PublicIPAllocationList struct {
	metav1.TypeMeta  `json:",inline"`
	*metav1.ListMeta `json:"metadata,omitempty"`
	Items            []PublicIPAllocation `json:"items"`
}
