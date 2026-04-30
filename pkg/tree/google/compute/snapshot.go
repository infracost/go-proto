package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Snapshot struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	StorageGB         value.Double `tree:"storage_gb"`
	SourceDisk        value.String `tree:"source_disk"`
}
