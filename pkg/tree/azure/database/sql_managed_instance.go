package database

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SQLManagedInstance struct {
	resource.Resource  `tree:"-"`
	SKU                value.String `tree:"sku"`
	LicenseType        value.Value[SQLLicenseType]      `tree:"license_type"`
	Cores              value.Int                        `tree:"cores"`
	StorageSizeGB      value.Int                        `tree:"storage_size_gb"`
	StorageAccountType value.Value[StorageAccountType]  `tree:"storage_account_type"`
}
