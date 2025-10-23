package k8s

import (
	"github.com/cybercoder/ik8s-ipam/pkg/api/v1alpha1/types"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

var publicIpPoolInformer cache.SharedIndexInformer

func CreatePublicIpPoolInformer() cache.SharedIndexInformer {
	if publicIpPoolInformer != nil {
		return publicIpPoolInformer
	}
	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(CreateDynamicClient(), 0, "", nil)
	publicIpPoolInformer = informerFactory.ForResource(types.PublicIPPoolGVR).Informer()
	return publicIpPoolInformer
}
