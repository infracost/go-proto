package mwaa

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Environment struct {
	resource.Resource `tree:"-"`
	EnvironmentClass  value.Value[EnvironmentClass] `tree:"environment_class"`
}

type EnvironmentClass uint32

const (
	EnvironmentClassUnknown   EnvironmentClass = iota
	EnvironmentClassMW1Small
	EnvironmentClassMW1Medium
	EnvironmentClassMW1Large
	EnvironmentClassMW1XLarge
	EnvironmentClassMW1XXLarge
)
