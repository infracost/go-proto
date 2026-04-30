package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FrontDoor struct {
	resource.Resource `tree:"-"`
	FrontendHosts     value.Int `tree:"frontend_hosts"`
	RoutingRules      value.Int `tree:"routing_rules"`
}
