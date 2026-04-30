package container

import (
	"github.com/infracost/go-proto/pkg/tree/google/storage"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Registry struct {
	resource.Resource `tree:"-"`
	Location          value.String                     `tree:"location"`
	StorageClass      value.Value[storage.StorageClass] `tree:"storage_class"`
}
