package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

type Image struct {
	resource.Resource `tree:"-"`
}
