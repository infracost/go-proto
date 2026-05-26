package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Disk struct {
	resource.Resource `tree:"-"`
	Name              value.String           `tree:"name"`
	SelfLink          value.String           `tree:"self_link"`
	Type              value.Value[DiskType]  `tree:"type"`
	StorageGB         value.Double           `tree:"storage_gb"`
	IOPS              value.Int              `tree:"iops"`
	InstanceCount     value.Int              `tree:"instance_count"`
	SourceImage       value.String           `tree:"source_image"`
	SourceSnapshot    value.String           `tree:"source_snapshot"`
	IsAttached        value.Bool             `tree:"is_attached"`
}

type DiskType uint32

const (
	DiskTypeUnknown            DiskType = iota
	DiskTypePDStandard
	DiskTypePDBalanced
	DiskTypePDSSD
	DiskTypePDExtreme
	DiskTypeHyperdiskBalanced
	DiskTypeHyperdiskExtreme
	DiskTypeHyperdiskThroughput
)
