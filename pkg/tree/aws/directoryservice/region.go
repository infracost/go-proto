package directoryservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Region is a replica region of a Microsoft AD (Enterprise) directory.
type Region struct {
	resource.Resource `tree:"-"`
	RegionName        value.String `tree:"region_name"`
	DomainControllers value.Int    `tree:"domain_controllers"`
}

type SharedDirectory struct {
	resource.Resource `tree:"-"`
	DirectoryID       value.String `tree:"directory_id"`

	Relationships SharedDirectoryRelationships `tree:"-"`
}

type SharedDirectoryRelationships struct {
	Directory *Directory
}
