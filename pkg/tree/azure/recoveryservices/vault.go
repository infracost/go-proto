package recoveryservices

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Vault struct {
	resource.Resource        `tree:"-"`
	Name                     value.String `tree:"name"`
	StorageModeType          value.Value[StorageModeType] `tree:"storage_mode_type"`
	CrossRegionRestoreEnabled value.Bool  `tree:"cross_region_restore_enabled"`

	Relationships VaultRelationships `tree:"-"`
}

type VaultRelationships struct {
	ProtectedVMs []*BackupProtectedVM
}
