package container

import (
	"github.com/infracost/go-proto/pkg/tree/google/compute"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type NodePool struct {
	resource.Resource  `tree:"-"`
	Name               value.String                      `tree:"name"`
	Zones              value.Int                         `tree:"zones"`
	CountPerZone       value.Int                         `tree:"count_per_zone"`
	MachineType        value.String                      `tree:"machine_type"`
	PurchaseOption     value.Value[compute.PurchaseOption] `tree:"purchase_option"`
	DiskType           value.Value[compute.DiskType]      `tree:"disk_type"`
	DiskSizeGB         value.Double                      `tree:"disk_size_gb"`
	LocalSSDCount      value.Int                         `tree:"local_ssd_count"`
	GuestAccelerators  []compute.GuestAccelerator        `tree:"guest_accelerators"`
	Inline             value.Bool                        `tree:"inline"`
	ClusterID          value.String                      `tree:"cluster_id"`

	Relationships NodePoolRelationships `tree:"-"`
}

type NodePoolRelationships struct {
	Cluster *Cluster
}
