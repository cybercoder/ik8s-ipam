package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cybercoder/ik8s-ipam/pkg/api/v1alpha1/types"
	"github.com/cybercoder/ik8s-ipam/pkg/k8s"
	"github.com/cybercoder/ik8s-ipam/pkg/netutils"
	"github.com/google/uuid"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
)

type IpamService struct {
	dynamicK8sClient *dynamic.DynamicClient
}

func NewIpamService() *IpamService {
	return &IpamService{
		dynamicK8sClient: k8s.CreateDynamicClient(),
	}
}

func (s *IpamService) FindOrCreateIPAssignment(r types.IpAssignmentRequestBody) (*types.PublicIPAllocation, error) {
	ctx := context.Background()
	list, err := s.dynamicK8sClient.Resource(types.PublicIPAllocationGVR).List(ctx, v1.ListOptions{
		FieldSelector: "spec.containerInterface=" + r.ContainerInterface + ",spec.resourceNamespace=" + r.Namespace + ",spec.resourceName=" + r.Name,
	})
	if err != nil {
		return nil, err
	}
	var alloc types.PublicIPAllocation
	if len(list.Items) < 1 {
		return s.createIPAssignment(r.ContainerInterface, r.Namespace, r.Name, r.IpFamily)
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.Items[0].Object, &alloc)
	if err != nil {
		return nil, fmt.Errorf("convert to PublicIPAllocation: %w", err)
	}

	return &alloc, nil

}

func (s *IpamService) createIPAssignment(containerInterface, namespace, name, ipFamily string) (*types.PublicIPAllocation, error) {
	ctx := context.Background()

	ipPool, err := s.findEmptyPublicIpPool(ipFamily)
	if err != nil {
		return nil, err
	}

	var idx uint64

	if len(ipPool.Status.ReleasedIndexes) > 0 {
		idx = ipPool.Status.ReleasedIndexes[0]
		ipPool.Status.ReleasedIndexes = ipPool.Status.ReleasedIndexes[1:]
		ipPool.Status.FreeIPs--
		ipPool.Status.AllocatedIPs++
	} else {
		idx = ipPool.Status.NextIndex
		ipPool.Status.FreeIPs--
		ipPool.Status.AllocatedIPs++
		if ipPool.Status.FreeIPs > 0 {
			ipPool.Status.NextIndex++
		}
	}

	ip, err := netutils.PickIPFromCIDRindex(ipPool.Spec.CIDR, idx)
	if err != nil {
		return nil, err
	}
	macAddress := netutils.GenerateStableMAC(containerInterface, namespace, name)
	alloc := types.PublicIPAllocation{
		ObjectMeta: &v1.ObjectMeta{
			Name: uuid.NewString(),
		},
		TypeMeta: v1.TypeMeta{
			Kind:       "PublicIPAllocation",
			APIVersion: "v1alpha1",
		},
		Spec: types.PublicIPAllocationSpec{
			MacAddress:         macAddress,
			ContainerInterface: containerInterface,
			IpFamily:           ipFamily,
			Address:            ip,
			ResourceKind:       "VirtualMachine",
			ResourceNamespace:  namespace,
			ResourceName:       name,
		},
		Status: types.PublicIPAllocationStatus{
			AllocatedAt: v1.NewTime(time.Now()),
			ReleasedAt:  nil,
		},
	}

	poolUnstruct, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ipPool)
	_, err = s.dynamicK8sClient.Resource(types.PublicIPPoolGVR).UpdateStatus(ctx, &unstructured.Unstructured{Object: poolUnstruct}, v1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	allocUnstruct, err := runtime.DefaultUnstructuredConverter.ToUnstructured(alloc)
	if err != nil {
		return nil, err
	}
	result, err := s.dynamicK8sClient.Resource(types.PublicIPAllocationGVR).Create(ctx, &unstructured.Unstructured{Object: allocUnstruct}, v1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	var allocResult types.PublicIPAllocation
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(result.Object, &allocResult)
	if err != nil {
		return nil, err
	}
	return &allocResult, nil
}

func (s *IpamService) findEmptyPublicIpPool(ipFamily string) (*types.PublicIPPool, error) {
	ctx := context.Background()
	list, err := s.dynamicK8sClient.Resource(types.PublicIPPoolGVR).List(ctx, v1.ListOptions{
		FieldSelector: "spec.ipFamily=" + ipFamily,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range list.Items {
		var pool types.PublicIPPool
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &pool); err != nil {
			return nil, err
		}
		if pool.Status.FreeIPs > 0 {
			return &pool, nil
		}
	}
	return nil, fmt.Errorf("no free %s pool found", ipFamily)
}
