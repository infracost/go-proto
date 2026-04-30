package compute

import "github.com/infracost/go-proto/pkg/tree/resource"

type HAVPNGateway struct {
	resource.Resource `tree:"-"`
}

type VPNTunnel struct {
	resource.Resource `tree:"-"`
}

type VPNGateway struct {
	resource.Resource `tree:"-"`
}
