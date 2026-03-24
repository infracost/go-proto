package apigatewayv2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type API struct {
	resource.Resource `tree:"-"`
	ProtocolType      value.String `tree:"protocol_type"`
}
