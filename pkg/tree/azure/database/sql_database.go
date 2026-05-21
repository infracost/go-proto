package database

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SQLDatabase struct {
	resource.Resource `tree:"-"`
	SKU               value.String                   `tree:"sku"`
	LicenseType       value.Value[SQLLicenseType]    `tree:"license_type"`
	Tier              value.String                   `tree:"tier"`
	Family            value.String                   `tree:"family"`
	Cores             value.Int                      `tree:"cores"`
	MaxSizeGB         value.Double                   `tree:"max_size_gb"`
	ReadReplicaCount  value.Int                      `tree:"read_replica_count"`
	ZoneRedundant     value.Bool                     `tree:"zone_redundant"`
	BackupStorageType value.Value[BackupStorageType] `tree:"backup_storage_type"`
	AutoPauseDelay    value.Int                      `tree:"auto_pause_delay"`
	MinCapacity       value.Double                   `tree:"min_capacity"`
}
