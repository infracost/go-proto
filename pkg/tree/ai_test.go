package tree

import (
	"testing"

	"github.com/infracost/go-proto/pkg/flag"
	"github.com/infracost/go-proto/pkg/tree/ai"
	"github.com/infracost/go-proto/pkg/tree/ai/code"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/infracost/proto/gen/go/infracost/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAIRoundTrip(t *testing.T) {
	configSource := &parser.SourceRange{Filename: "app/settings.py", StartLine: 12, EndLine: 12}

	original := &Tree{
		AI: ai.AI{
			Code: code.Code{
				Clients: []code.Client{
					{
						Resource: resource.Resource{ID: "client-1", Region: "us-east-1"},
						SDK:      value.New("boto3", 0, "", nil),
						Channel:  value.New(code.ChannelBedrock, 0, "", nil),
						Endpoint: value.New("https://bedrock-runtime.us-east-1.amazonaws.com", 0, "", nil),
					},
					{
						Resource: resource.Resource{ID: "client-2"},
						SDK:      value.New("openai", 0, "", nil),
						Channel:  value.New(code.ChannelProxy, 0, "", nil),
						Gateway:  value.New("litellm", 0, "", nil),
						Endpoint: value.New("http://litellm.internal:4000", 0, "", nil),
					},
				},
				ModelConfigs: []code.ModelConfig{
					{
						Resource: resource.Resource{ID: "config-1"},
						Kind:     value.New("env", 0, "", nil),
						Key:      value.New("CHAT_MODEL", 0, "", nil),
						Value:    value.New("claude-sonnet-4-5", flag.EnvVar, "", configSource),
					},
				},
				ModelCalls: []code.ModelCall{
					{
						// Resolved from an env var default set in another file.
						Resource:    resource.Resource{ID: "call-1"},
						SDK:         value.New("anthropic", 0, "", nil),
						Operation:   value.New("messages.create", 0, "", nil),
						Language:    value.New("python", 0, "", nil),
						ModelID:     value.New("claude-sonnet-4-5", flag.EnvVar, "", configSource),
						ModelSource: value.New("env:CHAT_MODEL", 0, "", nil),
						MaxTokens:   value.New[int64](1024, 0, "max_tokens", nil),
						// Explicitly disabled, as opposed to not set.
						PromptCaching: value.New(false, 0, "", nil),
						CacheTTL:      value.New("1h", 0, "", nil),
						Streaming:     value.New(true, 0, "", nil),
						Instrumentation: *value.NewList([]value.String{
							value.New("langfuse", 0, "", nil),
							value.New("opentelemetry", 0, "", nil),
						}, 0, "", nil),
						ClientID:      value.New("client-1", 0, "", nil),
						ModelConfigID: value.New("config-1", 0, "", nil),
					},
					{
						// Bedrock invoke_model: the model can't be resolved and
						// the settings sit inside a serialized body, so they
						// are all unknown.
						Resource:        resource.Resource{ID: "call-2"},
						SDK:             value.New("boto3", 0, "", nil),
						Operation:       value.New("invoke_model", 0, "", nil),
						ModelID:         value.New("", flag.Synthetic, "", nil),
						ModelSource:     value.New("unresolved", 0, "", nil),
						ModelExpression: value.New("settings.model", 0, "", nil),
						MaxTokens:       value.New[int64](0, flag.Synthetic, "", nil),
						PromptCaching:   value.New(false, flag.Synthetic, "", nil),
						ClientID:        value.New("client-1", 0, "", nil),
					},
					{
						// Through the proxy, as a batch job with thinking on.
						Resource:    resource.Resource{ID: "call-3"},
						SDK:         value.New("openai", 0, "", nil),
						Operation:   value.New("batches.create", 0, "", nil),
						ModelID:     value.New("gpt-5", 0, "", nil),
						ModelSource: value.New("literal", 0, "", nil),
						Batch:       value.New(true, 0, "", nil),
						Thinking:    value.New(true, 0, "", nil),
						ClientID:    value.New("client-2", 0, "", nil),
					},
				},
			},
		},
	}

	proto, err := original.ToProto()
	require.NoError(t, err)

	result, err := FromProto(proto)
	require.NoError(t, err)
	result.PostProcess()

	clients := result.AI.Code.Clients
	require.Len(t, clients, 2)
	assert.Equal(t, "us-east-1", clients[0].Region)
	assert.Equal(t, "boto3", clients[0].SDK.Value())
	assert.Equal(t, code.ChannelBedrock, clients[0].Channel.Value())
	assert.Empty(t, clients[0].Gateway.Value())
	assert.Equal(t, code.ChannelProxy, clients[1].Channel.Value())
	assert.Equal(t, "litellm", clients[1].Gateway.Value())
	assert.Equal(t, "http://litellm.internal:4000", clients[1].Endpoint.Value())

	configs := result.AI.Code.ModelConfigs
	require.Len(t, configs, 1)
	assert.Equal(t, "CHAT_MODEL", configs[0].Key.Value())
	assert.Equal(t, "claude-sonnet-4-5", configs[0].Value.Value())

	calls := result.AI.Code.ModelCalls
	require.Len(t, calls, 3)

	resolved := calls[0]
	assert.Equal(t, "claude-sonnet-4-5", resolved.ModelID.Value())
	assert.True(t, resolved.ModelID.Flags()&flag.EnvVar != 0)
	assert.False(t, resolved.ModelID.IsSynthetic())
	require.NotNil(t, resolved.ModelID.Source())
	assert.Equal(t, "app/settings.py", resolved.ModelID.Source().Filename)
	assert.Equal(t, int64(12), resolved.ModelID.Source().StartLine)
	assert.Equal(t, int64(1024), resolved.MaxTokens.Value())
	assert.False(t, resolved.MaxTokens.IsSynthetic())
	assert.False(t, resolved.PromptCaching.Value())
	assert.False(t, resolved.PromptCaching.IsSynthetic())
	assert.Equal(t, "1h", resolved.CacheTTL.Value())
	assert.True(t, resolved.Streaming.Value())
	// Not set: zero with no flags.
	assert.False(t, resolved.Batch.Value())
	assert.False(t, resolved.Batch.IsSynthetic())
	items := resolved.Instrumentation.Items()
	require.Len(t, items, 2)
	assert.Equal(t, "langfuse", items[0].Value())
	assert.Equal(t, "opentelemetry", items[1].Value())
	assert.Same(t, &result.AI.Code.Clients[0], resolved.Relationships.Client)
	assert.Same(t, &result.AI.Code.ModelConfigs[0], resolved.Relationships.ModelConfig)

	unresolved := calls[1]
	assert.Empty(t, unresolved.ModelID.Value())
	assert.True(t, unresolved.ModelID.IsSynthetic())
	assert.Equal(t, "unresolved", unresolved.ModelSource.Value())
	assert.Equal(t, "settings.model", unresolved.ModelExpression.Value())
	assert.Equal(t, int64(0), unresolved.MaxTokens.Value())
	assert.True(t, unresolved.MaxTokens.IsSynthetic())
	assert.True(t, unresolved.PromptCaching.IsSynthetic())
	assert.Same(t, &result.AI.Code.Clients[0], unresolved.Relationships.Client)
	assert.Nil(t, unresolved.Relationships.ModelConfig)

	proxied := calls[2]
	assert.True(t, proxied.Batch.Value())
	assert.True(t, proxied.Thinking.Value())
	assert.Same(t, &result.AI.Code.Clients[1], proxied.Relationships.Client)

	assert.Len(t, result.ToResources(false), 6)
}
