package lambda

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Function struct {
	resource.Resource `tree:"-"`
	Name              value.String       `tree:"name"`
	Architectures     value.List[string] `tree:"architectures"`
	MemorySizeMB      value.Int          `tree:"memory_size"`
	StorageSizeMB     value.Int          `tree:"storage_size"`
}
