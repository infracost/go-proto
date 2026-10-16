package msk

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

type ServerlessCluster struct {
	resource.Resource `tree:"-"`
}

type Replicator struct {
	resource.Resource `tree:"-"`
}

type VPCConnection struct {
	resource.Resource `tree:"-"`
}
