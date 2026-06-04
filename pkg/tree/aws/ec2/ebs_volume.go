package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EBSVolume struct {
	resource.Resource `tree:"-"`
	Type              value.Value[EBSVolumeType] `tree:"type"`
	IOPS              value.Int                  `tree:"iops"`
	SizeGB            value.Int                  `tree:"size"`
	ThroughputMiBperS value.Int                  `tree:"throughput"`
	MultiAttach       value.Bool                 `tree:"multi_attach"`
	Encrypted         value.Bool                 `tree:"encrypted"`
}

type EBSVolumeType uint32

const (
	EBSVolumeTypeUnknown  EBSVolumeType = iota
	EBSVolumeTypeGP2                    // General Purpose SSD (older gen)
	EBSVolumeTypeGP3                    // General Purpose SSD (current gen)
	EBSVolumeTypeIO1                    // Provisioned IOPS SSD (older gen)
	EBSVolumeTypeIO2                    // Provisioned IOPS SSD (current gen, includes Block Express)
	EBSVolumeTypeST1                    // Throughput Optimized HDD
	EBSVolumeTypeSC1                    // Cold HDD
	EBSVolumeTypeStandard               // Magnetic (previous generation)
)
