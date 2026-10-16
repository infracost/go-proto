package verifiedaccess

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Endpoint struct {
	resource.Resource `tree:"-"`
	// NonHTTP is true for TCP endpoints (rds, cidr, or tcp load-balancer/network-interface).
	NonHTTP value.Bool `tree:"non_http"`
}
