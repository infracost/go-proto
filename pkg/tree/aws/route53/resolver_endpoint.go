package route53

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ResolverEndpoint struct {
	resource.Resource  `tree:"-"`
	ResolverEndpoints  value.Int `tree:"resolver_endpoints"`
}
