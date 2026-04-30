package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Address struct {
	resource.Resource `tree:"-"`
	AddressType       value.Value[AddressType] `tree:"address_type"`
	Address           value.String             `tree:"address"`
	Purpose           value.Value[AddressPurpose] `tree:"purpose"`

	Relationships AddressRelationships `tree:"-"`
}

type AddressRelationships struct {
	Instance *Instance
}

type AddressPurpose uint32

const (
	AddressPurposeUnknown                AddressPurpose = iota
	AddressPurposeGCEEndpoint
	AddressPurposeDNSResolver
	AddressPurposeVPCPeering
	AddressPurposeIPsecInterconnect
	AddressPurposeSharedLoadbalancerVIP
	AddressPurposePrivateServiceConnect
)

type AddressType uint32

const (
	AddressTypeUnknown  AddressType = iota
	AddressTypeInternal
	AddressTypeExternal
)
