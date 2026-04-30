package stepfunctions

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type StateMachine struct {
	resource.Resource `tree:"-"`
	Type              value.Value[StateMachineType] `tree:"type"`
}

type StateMachineType uint32

const (
	StateMachineTypeUnknown  StateMachineType = iota
	StateMachineTypeStandard
	StateMachineTypeExpress
)
