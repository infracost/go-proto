package efs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// ReplicationConfiguration creates a destination file system that is billed as EFS storage.
type ReplicationConfiguration struct {
	resource.Resource `tree:"-"`
	// DestinationRegion is empty when the destination is in the same region.
	DestinationRegion   value.String `tree:"destination_region"`
	DestinationOneZone  value.Bool   `tree:"destination_one_zone"`
	ExistingDestination value.Bool   `tree:"existing_destination"`
}
