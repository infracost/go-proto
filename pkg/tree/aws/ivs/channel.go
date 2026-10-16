package ivs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Channel struct {
	resource.Resource `tree:"-"`
	Type              value.Value[ChannelType] `tree:"type"`
}

type ChannelType uint32

const (
	ChannelTypeUnknown ChannelType = iota
	ChannelTypeStandard
	ChannelTypeBasic
	ChannelTypeAdvancedSD
	ChannelTypeAdvancedHD
)
