package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPFamily represents IPv4 or IPv6
type IPFamily string

const (
	IPv4 IPFamily = "IPv4"
	IPv6 IPFamily = "IPv6"
)

// PublicIPPoolSpec defines the desired state of the pool
type PublicIPPoolSpec struct {
	IPFamily IPFamily `json:"ipFamily"` // IPv4 or IPv6
	CIDR     string   `json:"cidr"`     // Subnet in CIDR notation
	Gateway  string   `json:"gateway"`  // Optional gateway (pointer makes omitempty effective)
}

// PublicIPPoolStatus defines the observed state
type PublicIPPoolStatus struct {
	TotalIPs        uint64   `json:"totalIPs,omitempty"`
	AllocatedIPs    uint64   `json:"allocatedIPs,omitempty"`
	FreeIPs         uint64   `json:"freeIPs,omitempty"`
	NextIndex       uint64   `json:"nextIndex"`
	ReleasedIndexes []uint64 `json:"releasedIndexes"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ppool

// PublicIPPool is the Schema for PublicIPPools
type PublicIPPool struct {
	metav1.TypeMeta    `json:",inline"`
	*metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PublicIPPoolSpec   `json:"spec"`
	Status PublicIPPoolStatus `json:"status"`
}

// +kubebuilder:object:root=true

// PublicIPPoolList contains a list of PublicIPPool
type PublicIPPoolList struct {
	metav1.TypeMeta  `json:",inline"`
	*metav1.ListMeta `json:"metadata,omitempty"`
	Items            []PublicIPPool `json:"items"`
}
