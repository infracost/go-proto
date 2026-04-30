package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Snapshot struct {
	resource.Resource `tree:"-"`
	DiskSizeGB        value.Int    `tree:"disk_size_gb"`
	SourceDiskID      value.String `tree:"source_disk_id"`
}
