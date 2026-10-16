package redshiftserverless

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Workgroup struct {
	resource.Resource `tree:"-"`
	BaseCapacity      value.Int `tree:"base_capacity"`
}

type Namespace struct {
	resource.Resource `tree:"-"`
}

type Snapshot struct {
	resource.Resource `tree:"-"`
}
