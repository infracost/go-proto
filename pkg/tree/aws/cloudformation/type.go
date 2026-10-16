package cloudformation

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Type is a private registry extension (resource or hook).
type Type struct {
	resource.Resource `tree:"-"`
	Kind              value.Value[TypeKind] `tree:"kind"`
}

type TypeKind uint32

const (
	TypeKindUnknown TypeKind = iota
	TypeKindResource
	TypeKindHook
)
