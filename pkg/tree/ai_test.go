package tree

import (
	"testing"

	"github.com/infracost/go-proto/pkg/flag"
	"github.com/infracost/go-proto/pkg/tree/ai"
	"github.com/infracost/go-proto/pkg/tree/ai/code"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/infracost/proto/gen/go/infracost/parser"
	prototree "github.com/infracost/proto/gen/go/infracost/tree"
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
						Channel:  value.New(code.ChannelProxy, 0, "", nil),
						Provider: value.New("bedrock", 0, "", nil),
						// Still set for readers that predate Provider.
						Gateway:  value.New("bedrock", 0, "", nil),
						Endpoint: value.New("https://bedrock-runtime.us-east-1.amazonaws.com", 0, "", nil),
						// Formed from the client's region.
						EndpointSource: value.New("built", 0, "", nil),
					},
					{
						Resource:       resource.Resource{ID: "client-2"},
						SDK:            value.New("openai", 0, "", nil),
						Channel:        value.New(code.ChannelProxy, 0, "", nil),
						Provider:       value.New("litellm", 0, "", nil),
						Endpoint:       value.New("http://litellm.internal:4000", 0, "", nil),
						EndpointSource: value.New("env:LITELLM_BASE_URL", 0, "", nil),
					},
					{
						// Claude Platform on AWS: scoped to a workspace read
						// from an env var.
						Resource:  resource.Resource{ID: "client-3", Region: "us-west-2"},
						SDK:       value.New("anthropic", 0, "", nil),
						Channel:   value.New(code.ChannelProxy, 0, "", nil),
						Provider:  value.New("anthropic-aws", 0, "", nil),
						Workspace: value.New("wrkspc_01", flag.EnvVar, "", nil),
					},
					{
						// OpenAI's SDK pointed at Groq's API: direct, but
						// Groq runs the model.
						Resource: resource.Resource{ID: "client-4"},
						SDK:      value.New("openai", 0, "", nil),
						Channel:  value.New(code.ChannelDirect, 0, "", nil),
						Provider: value.New("groq", 0, "", nil),
						Endpoint: value.New("https://api.groq.com/openai/v1", 0, "", nil),
					},
					{
						// A local Ollama server.
						Resource: resource.Resource{ID: "client-5"},
						SDK:      value.New("ollama", 0, "", nil),
						Channel:  value.New(code.ChannelDirect, 0, "", nil),
						Provider: value.New("self-hosted", 0, "", nil),
					},
					{
						// From a parser that predates Provider: a proxy named
						// only in Gateway.
						Resource: resource.Resource{ID: "client-6"},
						SDK:      value.New("openai", 0, "", nil),
						Channel:  value.New(code.ChannelProxy, 0, "", nil),
						Gateway:  value.New("openrouter", 0, "", nil),
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
					{
						// A PydanticAI model class: a bare ID, which a
						// colon doesn't split into a provider and a model.
						Resource:    resource.Resource{ID: "call-4"},
						SDK:         value.New("pydantic-ai", 0, "", nil),
						Operation:   value.New("Agent", 0, "", nil),
						ModelID:     value.New("mistral:7b", 0, "", nil),
						ModelSource: value.New("literal", 0, "", nil),
						ModelClass:  value.New("OllamaModel", 0, "", nil),
					},
				},
			},
		},
	}

	// A writer runs PostProcess before encoding, as the parsers do.
	original.PostProcess()
	proto, err := original.ToProto()
	require.NoError(t, err)

	// On the wire, a proxy carries gateway as well as provider, whichever one
	// the writer set, for readers that predate Provider. A direct client
	// carries no gateway.
	assert.Equal(t, "litellm", wireString(t, proto, "client-2", "gateway"), "a proxy set only in Provider still writes gateway")
	assert.Equal(t, "litellm", wireString(t, proto, "client-2", "provider"))
	assert.Equal(t, "openrouter", wireString(t, proto, "client-6", "provider"), "a proxy set only in Gateway writes provider too")
	assert.Empty(t, wireString(t, proto, "client-4", "gateway"), "a direct client has no gateway")

	result, err := FromProto(proto)
	require.NoError(t, err)
	result.PostProcess()

	clients := result.AI.Code.Clients
	require.Len(t, clients, 6)
	assert.Equal(t, "us-east-1", clients[0].Region)
	assert.Equal(t, "boto3", clients[0].SDK.Value())
	assert.Equal(t, code.ChannelProxy, clients[0].Channel.Value())
	assert.Equal(t, "bedrock", clients[0].Provider.Value())
	assert.Equal(t, "bedrock", clients[0].Gateway.Value()) //nolint:staticcheck // Gateway still round-trips for older readers.
	assert.Equal(t, "built", clients[0].EndpointSource.Value())
	assert.Empty(t, clients[0].Workspace.Value())
	assert.Equal(t, code.ChannelProxy, clients[1].Channel.Value())
	assert.Equal(t, "litellm", clients[1].Provider.Value())
	assert.Equal(t, "http://litellm.internal:4000", clients[1].Endpoint.Value())
	assert.Equal(t, "env:LITELLM_BASE_URL", clients[1].EndpointSource.Value())
	assert.Equal(t, "anthropic-aws", clients[2].Provider.Value())
	assert.Equal(t, "wrkspc_01", clients[2].Workspace.Value())
	assert.True(t, clients[2].Workspace.Flags()&flag.EnvVar != 0)
	assert.Equal(t, code.ChannelDirect, clients[3].Channel.Value())
	assert.Equal(t, "groq", clients[3].Provider.Value())
	assert.Equal(t, "self-hosted", clients[4].Provider.Value())
	assert.Equal(t, "openrouter", clients[5].Provider.Value(), "an older tree's Gateway fills Provider")

	configs := result.AI.Code.ModelConfigs
	require.Len(t, configs, 1)
	assert.Equal(t, "CHAT_MODEL", configs[0].Key.Value())
	assert.Equal(t, "claude-sonnet-4-5", configs[0].Value.Value())

	calls := result.AI.Code.ModelCalls
	require.Len(t, calls, 4)

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
	assert.Empty(t, proxied.ModelClass.Value(), "a string model has no class")

	classed := calls[3]
	assert.Equal(t, "mistral:7b", classed.ModelID.Value())
	assert.Equal(t, "OllamaModel", classed.ModelClass.Value())
	assert.Nil(t, classed.Relationships.Client)

	assert.Len(t, result.ToResources(false), 11)
}

// wireString returns the string value of the entry key in the encoded
// attributes of the resource id, or "" when it has none.
func wireString(t *testing.T, tree *prototree.Tree, id, key string) string {
	t.Helper()
	for _, provider := range tree.GetProviders() {
		for _, service := range provider.GetServices() {
			for _, res := range service.GetResources() {
				if res.GetId() == id {
					return res.GetAttributes().GetEntries()[key].GetStringValue()
				}
			}
		}
	}
	t.Fatalf("resource %s not found in the encoded tree", id)
	return ""
}
