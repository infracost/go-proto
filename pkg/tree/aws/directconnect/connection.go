package directconnect

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Connection struct {
	resource.Resource `tree:"-"`
	Bandwidth         value.Value[Bandwidth] `tree:"bandwidth"`
	Location          value.String           `tree:"location"`
}

type Bandwidth uint32

const (
	BandwidthUnknown Bandwidth = iota
	Bandwidth50Mbps
	Bandwidth100Mbps
	Bandwidth200Mbps
	Bandwidth300Mbps
	Bandwidth400Mbps
	Bandwidth500Mbps
	Bandwidth1Gbps
	Bandwidth2Gbps
	Bandwidth5Gbps
	Bandwidth10Gbps
	Bandwidth25Gbps
	Bandwidth100Gbps
	Bandwidth400Gbps
)
