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

	// Gateway names the middleman when Channel is ChannelProxy, e.g.
	// "bedrock", "vertex", "azure-openai" or "openrouter". Who runs the
	// model, who bills for it and how its models are named are facts about
	// the gateway.
	Gateway value.String `tree:"gateway"`

	// Provider names who runs the model when Channel is ChannelDirect, since
	// "direct" says only that there's no middleman. It is the model's vendor,
	// such as "openai", "anthropic" or "moonshotai"; another host serving it,
	// such as "groq" or "together"; or "self-hosted" for an inference server
	// the code's owner runs, such as Ollama or vLLM. Empty when it isn't
	// known, and for a proxy, which Gateway names.
	Provider value.String `tree:"provider"`

	// Endpoint is the base URL the client is configured with, when set. The
	// parser must strip any userinfo and query string before setting it: base
	// URLs can carry credentials (https://user:key@host/, ?api-key=...), and
	// this value is shown in the dashboard and in PR comments.
	Endpoint value.String `tree:"endpoint"`

	// EndpointSource is how Endpoint was resolved: "literal", "const",
	// "env:<NAME>", "built" (formed from the client's options, as an SDK
	// forms its URL from a region or a resource name) or "unresolved". It
	// keeps an env var's name when its value can't be resolved and Endpoint
	// is empty.
	EndpointSource value.String `tree:"endpoint_source"`

	// Workspace is the workspace the client is scoped to, for clients that
	// take one, such as Claude Platform on AWS or Claude on Google Cloud. A
	// value resolved from an env var carries flag.EnvVar.
	Workspace value.String `tree:"workspace"`
}
