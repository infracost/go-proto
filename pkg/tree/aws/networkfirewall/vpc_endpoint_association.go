package networkfirewall

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VPCEndpointAssociation struct {
	resource.Resource `tree:"-"`
	EndpointCount     value.Int `tree:"endpoint_count"`
}
