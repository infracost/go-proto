package cloudwatch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LogGroup struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	RetentionInDays   value.Int    `tree:"retention_in_days"`
}
