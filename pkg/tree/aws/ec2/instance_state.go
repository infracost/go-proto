package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type InstanceStateMapping struct {
	resource.Resource `tree:"-"`
	InstanceID        value.String               `tree:"instance_id"`
	State             value.Value[InstanceState] `tree:"state"`
}

type InstanceState uint32

const (
	InstanceStateUnknown InstanceState = iota
	InstanceStateRunning
	InstanceStateStopped
)
