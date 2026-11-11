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

	// Provider names who the client sends its requests to. Read it with
	// Channel:
	//
	//   - ChannelDirect: who runs the model, since "direct" says only that
	//     there's no middleman. It is the model's vendor, such as "openai",
	//     "anthropic" or "moonshotai"; another host serving it, such as "groq"
	//     or "together"; or "self-hosted" for an inference server the code's
	//     owner runs, such as Ollama or vLLM.
	//   - ChannelProxy: the middleman, such as "bedrock", "vertex",
	//     "azure-openai" or "openrouter". Who runs the model, who bills for it
	//     and how its models are named are facts about the middleman.
	//
	// Empty when it isn't known. It names who the client calls, not who made
	// the model: one client can call several vendors' models, and a call's
	// model says whose it is.
	Provider value.String `tree:"provider"`

	// Gateway names the middleman when Channel is ChannelProxy, the same
	// value as Provider. Parsers still set it for readers that predate
	// Provider.
	//
	// Deprecated: use Provider, which is set for every channel.
	Gateway value.String `tree:"gateway"`

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
