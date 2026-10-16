package kinesisvideo

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Stream struct {
	resource.Resource    `tree:"-"`
	DataRetentionInHours value.Int `tree:"data_retention_in_hours"`
}
