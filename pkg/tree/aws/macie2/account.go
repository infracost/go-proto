package macie2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Account struct {
	resource.Resource `tree:"-"`
	Enabled           value.Bool `tree:"enabled"`
}
