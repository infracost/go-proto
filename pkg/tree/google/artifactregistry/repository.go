package artifactregistry

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Repository struct {
	resource.Resource `tree:"-"`
	Location          value.String `tree:"location"`
}
