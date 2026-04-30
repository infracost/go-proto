package cloudfunctions

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Function struct {
	resource.Resource `tree:"-"`
	AvailableMemoryMB value.Int `tree:"available_memory_mb"`
}
