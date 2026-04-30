package cosmosdb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Account struct {
	resource.Resource                `tree:"-"`
	Name                             value.String      `tree:"name"`
	OfferType                        value.String      `tree:"offer_type"`
	Kind                             value.Value[CosmosDBKind]      `tree:"kind"`
	ConsistencyLevel                 value.Value[ConsistencyLevel]  `tree:"consistency_level"`
	GeoLocations                     []GeoLocation     `tree:"geo_locations"`
	Capabilities                     value.List[string] `tree:"capabilities"`
	EnableMultipleWriteLocations     value.Bool         `tree:"enable_multiple_write_locations"`
	EnableAutomaticFailover          value.Bool         `tree:"enable_automatic_failover"`
	AnalyticalStorageEnabled         value.Bool         `tree:"analytical_storage_enabled"`
	FreeTierEnabled                  value.Bool         `tree:"free_tier_enabled"`
	BackupType                       value.Value[BackupType]              `tree:"backup_type"`
	BackupStorageRedundancy          value.Value[BackupStorageRedundancy] `tree:"backup_storage_redundancy"`
	BackupIntervalInMinutes          value.Int          `tree:"backup_interval_in_minutes"`
	BackupRetentionInHours           value.Int          `tree:"backup_retention_in_hours"`
}

type GeoLocation struct {
	Location      value.String `tree:"location"`
	ZoneRedundant value.Bool   `tree:"zone_redundant"`
}
