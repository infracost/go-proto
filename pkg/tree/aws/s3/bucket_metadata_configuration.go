package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BucketMetadataConfiguration struct {
	resource.Resource `tree:"-"`
	InventoryEnabled  value.Bool `tree:"inventory_enabled"`
}
