package datastream

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

type Stream struct {
	resource.Resource `tree:"-"`
}
