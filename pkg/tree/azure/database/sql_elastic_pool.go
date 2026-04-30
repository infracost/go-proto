package database

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SQLElasticPool struct {
	resource.Resource `tree:"-"`
	SKU               value.String `tree:"sku"`
	Tier              value.String `tree:"tier"`
	Family            value.String `tree:"family"`
	Cores             value.Int    `tree:"cores"`
	DTUCapacity       value.Int    `tree:"dtu_capacity"`
	MaxSizeGB         value.Double `tree:"max_size_gb"`
	ZoneRedundant     value.Bool   `tree:"zone_redundant"`
	LicenseType       value.Value[SQLLicenseType] `tree:"license_type"`
}
