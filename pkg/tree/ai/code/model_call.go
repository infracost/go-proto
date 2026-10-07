package code

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// ModelCall is one call site that sends a request to a model. The base
// resource's source range points at the call.
//
// The settings below follow the set / not set / unknown convention described
// in the package doc.
type ModelCall struct {
	resource.Resource `tree:"-"`

	// SDK is the SDK the call is made with, e.g. "anthropic" or "openai".
	SDK value.String `tree:"sdk"`

	// Operation is the SDK method called, e.g. "messages.create" or
	// "invoke_model".
	Operation value.String `tree:"operation"`

	// Language is the source language of the call, e.g. "python".
	Language value.String `tree:"language"`

	// ModelID is the model the call uses. Its source range points at where
	// the model was set, which may be a different file from the call. A value
	// resolved from an env var or a config file carries flag.EnvVar or
	// flag.Config.
	//
	// A model the parser could not resolve is an empty value with
	// flag.Synthetic, and the expression it could not resolve is kept in
	// ModelExpression.
	ModelID value.String `tree:"model_id"`

	// ModelSource is how ModelID was resolved: "literal", "const",
	// "env:<NAME>", "config" or "unresolved". It is a summary for display and
	// grouping. Where it overlaps other fields, those are authoritative:
	// ModelID's flags (Synthetic, EnvVar, Config) say whether and how the
	// model was resolved, and the linked ModelConfig's Key names the env var
	// or config key.
	ModelSource value.String `tree:"model_source"`

	// ModelExpression is the raw source expression the model is passed as,
	// e.g. "settings.model". It is what to read when ModelID is unresolved.
	ModelExpression value.String `tree:"model_expression"`

	// MaxTokens is the output token limit, resolved to an integer.
	MaxTokens value.Int `tree:"max_tokens"`

	// PromptCaching is whether the call marks any of its prompt for caching.
	PromptCaching value.Bool `tree:"prompt_caching"`

	// CacheTTL is the cache lifetime the call asks for, e.g. "5m" or "1h".
	CacheTTL value.String `tree:"cache_ttl"`

	// Streaming is whether the call streams its response.
	Streaming value.Bool `tree:"streaming"`

	// Batch is whether the call goes through a batch API.
	Batch value.Bool `tree:"batch"`

	// Thinking is whether the call enables extended thinking or reasoning.
	Thinking value.Bool `tree:"thinking"`

	// Instrumentation lists the observability tools wrapping the call, e.g.
	// "langfuse" or "opentelemetry".
	Instrumentation value.List[string] `tree:"instrumentation"`

	// ClientID is the ID of the Client the call is made through, if found.
	ClientID value.String `tree:"client_id"`

	// ModelConfigID is the ID of the ModelConfig ModelID was read from, if it
	// was set away from the call.
	ModelConfigID value.String `tree:"model_config_id"`

	Relationships ModelCallRelationships `tree:"-"`
}

type ModelCallRelationships struct {
	Client      *Client
	ModelConfig *ModelConfig
}
