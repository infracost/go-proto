package database

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MySQLFlexibleServer struct {
	resource.Resource         `tree:"-"`
	SKUName                   value.String `tree:"sku_name"`
	StorageSizeGB             value.Int    `tree:"storage_size_gb"`
	IOPS                      value.Int    `tree:"iops"`
	GeoRedundantBackupEnabled value.Bool   `tree:"geo_redundant_backup_enabled"`
}
