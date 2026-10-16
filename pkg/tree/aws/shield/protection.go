package shield

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Protection struct {
	resource.Resource `tree:"-"`
	ResourceType      value.Value[ProtectedResourceType] `tree:"resource_type"`
}

type ProtectedResourceType uint32

const (
	ProtectedResourceTypeUnknown ProtectedResourceType = iota
	ProtectedResourceTypeLoadBalancer
	ProtectedResourceTypeElasticIP
	ProtectedResourceTypeCloudFront
	ProtectedResourceTypeGlobalAccelerator
	ProtectedResourceTypeRoute53
)
