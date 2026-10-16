package emr

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ServerlessApplication struct {
	resource.Resource `tree:"-"`
	Architecture      value.Value[Architecture] `tree:"architecture"`
}

type Architecture uint32

const (
	ArchitectureUnknown Architecture = iota
	ArchitectureX86_64
	ArchitectureARM64
)
