package types

import "k8s.io/apimachinery/pkg/runtime/schema"

var PublicIPPoolGVR = schema.GroupVersionResource{
	Group:    "ovn.ik8s.ir",
	Version:  "v1alpha1",
	Resource: "publicippools",
}

var PublicIPAllocationGVR = schema.GroupVersionResource{
	Group:    "ovn.ik8s.ir",
	Version:  "v1alpha1",
	Resource: "publicipallocations",
}
