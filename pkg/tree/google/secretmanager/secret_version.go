package secretmanager

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SecretVersion struct {
	resource.Resource    `tree:"-"`
	ReplicationLocations value.Int    `tree:"replication_locations"`
	SecretRef            value.String `tree:"secret"`

	Relationships SecretVersionRelationships `tree:"-"`
}

type SecretVersionRelationships struct {
	Secret *Secret
}
