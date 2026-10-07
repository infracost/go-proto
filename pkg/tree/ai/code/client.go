package code

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Client is the construction of an SDK client, such as anthropic.Anthropic()
// or boto3.client("bedrock-runtime"). The region a cloud-routed client sends
// requests to is the base resource's Region.
type Client struct {
	resource.Resource `tree:"-"`

	// SDK is the SDK the client is constructed from, e.g. "anthropic" or
	// "boto3".
	SDK value.String `tree:"sdk"`

	// Channel is the route the client's requests take to the model.
	Channel value.Value[Channel] `tree:"channel"`

	// Gateway names the gateway or proxy when Channel is ChannelProxy, e.g.
	// "litellm".
	Gateway value.String `tree:"gateway"`

	// Endpoint is the base URL the client is configured with, when set. The
	// parser must strip any userinfo and query string before setting it: base
	// URLs can carry credentials (https://user:key@host/, ?api-key=...), and
	// this value is shown in the dashboard and in PR comments.
	Endpoint value.String `tree:"endpoint"`
}
