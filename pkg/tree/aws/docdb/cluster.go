package docdb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource     `tree:"-"`
	BackupRetentionPeriod value.Int    `tree:"backup_retention_period"`
	EngineVersion         value.String `tree:"engine_version"`
	StorageEncrypted      value.Bool   `tree:"storage_encrypted"`
}
