package redis

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Instance struct {
	resource.Resource `tree:"-"`
	Tier              value.Value[RedisTier] `tree:"tier"`
	MemorySizeGB      value.Double           `tree:"memory_size_gb"`
}

type RedisTier uint32

const (
	RedisTierUnknown    RedisTier = iota
	RedisTierBasic
	RedisTierStandardHA
)
