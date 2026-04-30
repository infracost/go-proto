package secretmanager

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Secret struct {
	resource.Resource    `tree:"-"`
	Name                 value.String `tree:"name"`
	ReplicationLocations value.Int    `tree:"replication_locations"`
}
