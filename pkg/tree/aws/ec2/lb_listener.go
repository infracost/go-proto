package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LBListener struct {
	resource.Resource `tree:"-"`
	LoadBalancerARN   value.String                  `tree:"load_balancer_arn"`
	Port              value.Int                     `tree:"port"`
	Protocol          value.Value[LBListenerProtocol] `tree:"protocol"`
	DefaultActions    []LBListenerDefaultAction     `tree:"default_action"`
}

type LBListenerDefaultAction struct {
	Type     value.Value[LBListenerActionType] `tree:"type"`
	Redirect *LBListenerRedirectAction          `tree:"redirect"`
}

type LBListenerActionType uint32

const (
	LBListenerActionTypeUnknown           LBListenerActionType = iota
	LBListenerActionTypeForward
	LBListenerActionTypeRedirect
	LBListenerActionTypeFixedResponse
	LBListenerActionTypeAuthenticateCognito
	LBListenerActionTypeAuthenticateOIDC
)

type LBListenerRedirectAction struct {
	Protocol   value.Value[LBListenerRedirectProtocol] `tree:"protocol"`
	Port       value.String                            `tree:"port"`
	StatusCode value.Value[LBListenerRedirectStatusCode] `tree:"status_code"`
}

type LBListenerProtocol uint32

const (
	LBListenerProtocolUnknown LBListenerProtocol = iota
	LBListenerProtocolHTTP
	LBListenerProtocolHTTPS
	LBListenerProtocolTCP
	LBListenerProtocolTLS
	LBListenerProtocolUDP
	LBListenerProtocolTCPUDP
)

type LBListenerRedirectStatusCode uint32

const (
	LBListenerRedirectStatusCodeUnknown LBListenerRedirectStatusCode = iota
	LBListenerRedirectStatusCodeHTTP301
	LBListenerRedirectStatusCodeHTTP302
)

type LBListenerRedirectProtocol uint32

const (
	LBListenerRedirectProtocolUnknown LBListenerRedirectProtocol = iota
	LBListenerRedirectProtocolHTTP
	LBListenerRedirectProtocolHTTPS
)
