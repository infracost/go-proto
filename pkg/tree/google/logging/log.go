package logging

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BucketConfig struct {
	resource.Resource `tree:"-"`
	Location          value.String `tree:"location"`
}
