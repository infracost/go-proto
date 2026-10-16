package fsx

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FileCache struct {
	resource.Resource `tree:"-"`
	StorageCapacityGB value.Int `tree:"storage_capacity_gb"`
}
