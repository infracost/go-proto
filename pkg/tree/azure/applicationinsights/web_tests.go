package applicationinsights

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type StandardWebTest struct {
	resource.Resource `tree:"-"`
	Enabled           value.Bool `tree:"enabled"`
	FrequencySeconds  value.Int  `tree:"frequency"`
}

type WebTest struct {
	resource.Resource `tree:"-"`
	Kind              value.String `tree:"kind"`
	Enabled           value.Bool   `tree:"enabled"`
}
