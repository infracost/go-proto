package applicationinsights

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Insights struct {
	resource.Resource `tree:"-"`
	RetentionInDays   value.Int `tree:"retention_in_days"`
}
