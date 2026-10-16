package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type StorageLensConfiguration struct {
	resource.Resource `tree:"-"`
	Enabled           value.Bool `tree:"enabled"`
	// AdvancedMetrics is true when any paid (advanced) metric group is enabled.
	AdvancedMetrics value.Bool `tree:"advanced_metrics"`
}
