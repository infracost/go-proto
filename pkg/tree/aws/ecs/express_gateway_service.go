package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ExpressGatewayService struct {
	resource.Resource `tree:"-"`
	CPU               value.Int `tree:"cpu"`
	MemoryMiB         value.Int `tree:"memory"`
	MinTaskCount      value.Int `tree:"min_task_count"`
}
