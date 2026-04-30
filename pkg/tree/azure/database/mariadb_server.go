package database

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MariaDBServer struct {
	resource.Resource         `tree:"-"`
	SKUName                   value.String `tree:"sku_name"`
	StorageMB                 value.Int    `tree:"storage_mb"`
	GeoRedundantBackupEnabled value.Bool   `tree:"geo_redundant_backup_enabled"`
}
