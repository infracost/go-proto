package storagegateway

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Gateway struct {
	resource.Resource `tree:"-"`
	GatewayType       value.Value[GatewayType] `tree:"gateway_type"`
}

type GatewayType uint32

const (
	GatewayTypeUnknown GatewayType = iota
	GatewayTypeStored
	GatewayTypeCached
	GatewayTypeVTL
	GatewayTypeVTLSnow
	GatewayTypeFileS3
	GatewayTypeFileFSxSMB
)
