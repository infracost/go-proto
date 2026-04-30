package spanner

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Instance struct {
	resource.Resource    `tree:"-"`
	Name                 value.String                   `tree:"name"`
	NumNodes             value.Int                      `tree:"num_nodes"`
	ProcessingUnits      value.Int                      `tree:"processing_units"`
	Edition              value.Value[SpannerEdition]    `tree:"edition"`
	Config               value.String                   `tree:"config"`
	InstanceType         value.Value[SpannerInstanceType] `tree:"instance_type"`
	HasManagedAutoscaler bool                           `tree:"-"`
}

type SpannerInstanceType uint32

const (
	SpannerInstanceTypeUnknown      SpannerInstanceType = iota
	SpannerInstanceTypeProvisioned
	SpannerInstanceTypeFreeInstance
)

type SpannerEdition uint32

const (
	SpannerEditionUnknown        SpannerEdition = iota
	SpannerEditionStandard
	SpannerEditionEnterprise
	SpannerEditionEnterprisePlus
)
