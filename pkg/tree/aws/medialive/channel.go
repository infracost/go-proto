package medialive

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Channel struct {
	resource.Resource `tree:"-"`
	ChannelClass      value.Value[ChannelClass]    `tree:"channel_class"`
	InputCodec        value.Value[InputCodec]      `tree:"input_codec"`
	InputResolution   value.Value[InputResolution] `tree:"input_resolution"`
	MaximumBitrate    value.Value[MaximumBitrate]  `tree:"maximum_bitrate"`
}

type ChannelClass uint32

const (
	ChannelClassUnknown ChannelClass = iota
	ChannelClassStandard
	ChannelClassSinglePipeline
)

type InputCodec uint32

const (
	InputCodecUnknown InputCodec = iota
	InputCodecAVC
	InputCodecHEVC
	InputCodecMPEG2
)

type InputResolution uint32

const (
	InputResolutionUnknown InputResolution = iota
	InputResolutionSD
	InputResolutionHD
	InputResolutionUHD
)

type MaximumBitrate uint32

const (
	MaximumBitrateUnknown MaximumBitrate = iota
	MaximumBitrateMax10Mbps
	MaximumBitrateMax20Mbps
	MaximumBitrateMax50Mbps
)
