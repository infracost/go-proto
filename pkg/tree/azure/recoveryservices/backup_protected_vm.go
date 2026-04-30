package recoveryservices

import (
	"github.com/infracost/go-proto/pkg/tree/azure/compute"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BackupProtectedVM struct {
	resource.Resource `tree:"-"`
	RecoveryVaultName value.String `tree:"recovery_vault_name"`
	SourceVMID        value.String `tree:"source_vm_id"`

	Relationships BackupProtectedVMRelationships `tree:"-"`
}

type BackupProtectedVMRelationships struct {
	SourceVM *compute.VirtualMachine
}
