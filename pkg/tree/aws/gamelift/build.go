package gamelift

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Build struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	OperatingSystem   value.String `tree:"operating_system"`
}
