package sql

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DatabaseInstance struct {
	resource.Resource `tree:"-"`
	DiskSizeGB        value.Int                        `tree:"disk_size_gb"`
	UseIPV4           value.Bool                       `tree:"use_ipv4"`
	HasReplica        value.Bool                       `tree:"has_replica"`
	Tier              value.String                     `tree:"tier"`
	Edition           value.Value[SQLEdition]          `tree:"edition"`
	AvailabilityType  value.Value[AvailabilityType]    `tree:"availability_type"`
	DatabaseVersion   value.String                     `tree:"database_version"`
	DiskType          value.Value[SQLDiskType]         `tree:"disk_type"`
}

type SQLEdition uint32

const (
	SQLEditionUnknown        SQLEdition = iota
	SQLEditionEnterprise
	SQLEditionEnterprisePlus
)

type AvailabilityType uint32

const (
	AvailabilityTypeUnknown  AvailabilityType = iota
	AvailabilityTypeZonal
	AvailabilityTypeRegional
)

type SQLDiskType uint32

const (
	SQLDiskTypeUnknown SQLDiskType = iota
	SQLDiskTypePDSSD
	SQLDiskTypePDHDD
)
