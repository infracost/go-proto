package network

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &Network{
		LoadBalancers: []LoadBalancer{
			{Resource: resource.Resource{ID: "lb-1"}},
		},
		LoadBalancerRules: []LoadBalancerRule{
			{
				Resource:       resource.Resource{ID: "lbr-1"},
				LoadBalancerID: value.New("lb-1", 0, "", nil),
			},
		},
		VirtualNetworkGateways: []VirtualNetworkGateway{
			{Resource: resource.Resource{ID: "vng-1"}},
		},
		VirtualNetworkGatewayConnections: []VirtualNetworkGatewayConnection{
			{
				Resource:                resource.Resource{ID: "vngc-1"},
				VirtualNetworkGatewayID: value.New("vng-1", 0, "", nil),
			},
		},
		VirtualNetworks: []VirtualNetwork{
			{
				Resource: resource.Resource{ID: "vnet-1"},
				Name:     value.New("vnet-1", 0, "", nil),
			},
		},
		VirtualNetworkPeerings: []VirtualNetworkPeering{
			{
				Resource:               resource.Resource{ID: "vnp-1"},
				RemoteVirtualNetworkID: value.New("vnet-1", 0, "", nil),
				VirtualNetworkName:     value.New("vnet-1", 0, "", nil),
			},
		},
		TrafficManagerProfiles: []TrafficManagerProfile{
			{Resource: resource.Resource{ID: "tmp-1"}},
		},
		TrafficManagerEndpoints: []TrafficManagerEndpoint{
			{
				Resource:  resource.Resource{ID: "tme-1"},
				ProfileID: value.New("tmp-1", 0, "", nil),
			},
		},
		CDNProfiles: []CDNProfile{
			{
				Resource: resource.Resource{ID: "cp-1"},
				Name:     value.New("cp-1", 0, "", nil),
			},
		},
		CDNEndpoints: []CDNEndpoint{
			{
				Resource:    resource.Resource{ID: "ce-1"},
				ProfileName: value.New("cp-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	lbRule := s.LoadBalancerRules[0].Relationships.LoadBalancer
	gwConn := s.VirtualNetworkGatewayConnections[0].Relationships.Gateway
	peerSrc := s.VirtualNetworkPeerings[0].Relationships.SourceVirtualNetwork
	peerRem := s.VirtualNetworkPeerings[0].Relationships.RemoteVirtualNetwork
	tmEnd := s.TrafficManagerEndpoints[0].Relationships.Profile
	cdnEnd := s.CDNEndpoints[0].Relationships.Profile

	s.PostProcess()
	assert.Equal(t, lbRule, s.LoadBalancerRules[0].Relationships.LoadBalancer)
	assert.Equal(t, gwConn, s.VirtualNetworkGatewayConnections[0].Relationships.Gateway)
	assert.Equal(t, peerSrc, s.VirtualNetworkPeerings[0].Relationships.SourceVirtualNetwork)
	assert.Equal(t, peerRem, s.VirtualNetworkPeerings[0].Relationships.RemoteVirtualNetwork)
	assert.Equal(t, tmEnd, s.TrafficManagerEndpoints[0].Relationships.Profile)
	assert.Equal(t, cdnEnd, s.CDNEndpoints[0].Relationships.Profile)

	assert.NotNil(t, lbRule)
	assert.NotNil(t, gwConn)
	assert.NotNil(t, peerSrc)
	assert.NotNil(t, peerRem)
	assert.NotNil(t, tmEnd)
	assert.NotNil(t, cdnEnd)
}
