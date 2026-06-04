package neptune

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource           `tree:"-"`
	Identifier                  value.String                    `tree:"identifier"`
	BackupRetentionPeriod       value.Int                       `tree:"backup_retention_period"`
	StorageType                 value.Value[NeptuneStorageType] `tree:"storage_type"`
	EnableCloudwatchLogsExports value.List[string]              `tree:"enable_cloudwatch_logs_exports"`
	StorageEncrypted            value.Bool                      `tree:"storage_encrypted"`
}

type NeptuneStorageType uint32

const (
	NeptuneStorageTypeUnknown  NeptuneStorageType = iota
	NeptuneStorageTypeStandard
	NeptuneStorageTypeIOPT1
)
