package apigatewayv2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type API struct {
	resource.Resource `tree:"-"`
	ProtocolType      value.Value[ProtocolType] `tree:"protocol_type"`
}

type ProtocolType uint32

const (
	ProtocolTypeUnknown   ProtocolType = iota
	ProtocolTypeHTTP
	ProtocolTypeWebSocket
)
