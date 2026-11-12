package code

import (
	"testing"

	"github.com/infracost/go-proto/pkg/flag"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeModelID_Prefixes(t *testing.T) {
	for _, tc := range []struct {
		sdk, id, normalized, route string
	}{
		// PydanticAI's and LangChain's provider:model.
		{"pydantic-ai", "openai:gpt-4o", "gpt-4o", "openai"},
		{"pydantic-ai", "anthropic:claude-sonnet-4-5", "claude-sonnet-4-5", "anthropic"},
		{"pydantic-ai", "bedrock:anthropic.claude-3-5-sonnet-20240620-v1:0", "anthropic.claude-3-5-sonnet-20240620-v1:0", "bedrock"},
		{"pydantic-ai", "google-gla:gemini-2.0-flash", "gemini-2.0-flash", "google-gla"},
		{"langchain", "bedrock_converse:anthropic.claude", "anthropic.claude", "bedrock_converse"},
		// LiteLLM's and instructor's provider/model.
		{"litellm", "bedrock/anthropic.claude-3-sonnet", "anthropic.claude-3-sonnet", "bedrock"},
		{"litellm", "azure/my-gpt4o-deployment", "my-gpt4o-deployment", "azure"},
		{"litellm", "openrouter/meta-llama/llama-3-70b", "meta-llama/llama-3-70b", "openrouter"},
		{"litellm", "fireworks_ai/accounts/fireworks/models/llama-v3-70b", "accounts/fireworks/models/llama-v3-70b", "fireworks_ai"},
		{"litellm", "ollama/llama3", "llama3", "ollama"},
		{"litellm", "ollama_chat/llama3", "llama3", "ollama_chat"},
		{"litellm", "hosted_vllm/meta-llama/Llama-3-8b", "meta-llama/Llama-3-8b", "hosted_vllm"},
		{"litellm", "sagemaker/my-endpoint", "my-endpoint", "sagemaker"},
		{"litellm", "replicate/meta/llama-3-70b", "meta/llama-3-70b", "replicate"},
		{"litellm", "watsonx/ibm/granite-13b-chat-v2", "ibm/granite-13b-chat-v2", "watsonx"},
		{"litellm", "nvidia_nim/meta/llama3-70b-instruct", "meta/llama3-70b-instruct", "nvidia_nim"},
		{"langchain", "ibm:ibm/granite-3-8b-instruct", "ibm/granite-3-8b-instruct", "ibm"},
		{"langchain", "nvidia:meta/llama3-70b-instruct", "meta/llama3-70b-instruct", "nvidia"},
		// LiteLLM's Bedrock route picks its API with a second segment.
		{"litellm", "bedrock/converse/us.anthropic.claude-3-5-sonnet-20240620-v1:0", "us.anthropic.claude-3-5-sonnet-20240620-v1:0", "bedrock"},
		{"litellm", "bedrock/invoke/anthropic.claude-3-5-sonnet-20240620-v1:0", "anthropic.claude-3-5-sonnet-20240620-v1:0", "bedrock"},
		// Only Bedrock's route has those, and only as a segment of their own.
		{"litellm", "bedrock/converse", "converse", "bedrock"},
		{"litellm", "openai/converse/gpt-4o", "converse/gpt-4o", "openai"},
		{"pydantic-ai", "bedrock:converse:anthropic.claude", "converse:anthropic.claude", "bedrock"},
		{"instructor", "openai/gpt-4o-mini", "gpt-4o-mini", "openai"},
		{"litellm", "bedrock/arn:aws:bedrock:us-east-1:123456789012:application-inference-profile/abc", "arn:aws:bedrock:us-east-1:123456789012:application-inference-profile/abc", "bedrock"},
		// A separator the SDK doesn't use is part of the ID.
		{"litellm", "anthropic.claude-3-5-sonnet-20240620-v1:0", "anthropic.claude-3-5-sonnet-20240620-v1:0", ""},
		{"pydantic-ai", "openai/gpt-4o", "openai/gpt-4o", ""},
		{"litellm", "meta-llama/Llama-3-70b-chat-hf", "meta-llama/Llama-3-70b-chat-hf", ""},
		// SDKs that take the model ID as is.
		{"openai", "mistral:7b", "mistral:7b", ""},
		{"openai", "anthropic/claude-3.5-sonnet", "anthropic/claude-3.5-sonnet", ""},
		{"together", "meta-llama/Llama-3-70b-chat-hf", "meta-llama/Llama-3-70b-chat-hf", ""},
		{"fireworks-ai", "accounts/fireworks/models/llama-v3-70b", "accounts/fireworks/models/llama-v3-70b", ""},
		{"boto3", "us.anthropic.claude-3-5-sonnet-20240620-v1:0", "us.anthropic.claude-3-5-sonnet-20240620-v1:0", ""},
		{"", "openai:gpt-4o", "openai:gpt-4o", ""},
		// Plain, empty and degenerate input.
		{"pydantic-ai", "gpt-4o", "gpt-4o", ""},
		{"pydantic-ai", "", "", ""},
		{"pydantic-ai", "openai:", "openai:", ""},
		{"pydantic-ai", ":gpt-4o", ":gpt-4o", ""},
		{"pydantic-ai", "llama3:8b", "llama3:8b", ""},
		// The string form, Agent("mistral:7b"), names PydanticAI's mistral
		// provider, so it is stripped. The model class form,
		// OllamaModel("mistral:7b"), is never normalized: see
		// TestModelCall_NormalizedModelID.
		{"pydantic-ai", "mistral:7b", "7b", "mistral"},
	} {
		t.Run(tc.sdk+" "+tc.id, func(t *testing.T) {
			normalized, route := NormalizeModelID(tc.sdk, tc.id)
			assert.Equal(t, tc.normalized, normalized)
			assert.Equal(t, tc.route, route)
		})
	}
}

func TestModelCall_NormalizedModelID(t *testing.T) {
	str := func(v string) value.String { return value.New(v, 0, "", nil) }
	for _, tc := range []struct {
		name              string
		call              ModelCall
		normalized, route string
	}{
		{"a string ID with its SDK's prefix", ModelCall{SDK: str("pydantic-ai"), ModelID: str("openai:gpt-4o")}, "gpt-4o", "openai"},
		{"an SDK that takes the ID as is", ModelCall{SDK: str("openai"), ModelID: str("mistral:7b")}, "mistral:7b", ""},
		// The model class names the ID as its provider does: "mistral:7b" is
		// Ollama's tag, not PydanticAI's mistral provider.
		{"a model class's ID is never normalized", ModelCall{SDK: str("pydantic-ai"), ModelClass: str("OllamaModel"), ModelID: str("mistral:7b")}, "mistral:7b", ""},
		{"an unresolved model", ModelCall{SDK: str("pydantic-ai"), ModelID: value.New("", flag.Synthetic, "", nil)}, "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			normalized, route := tc.call.NormalizedModelID()
			assert.Equal(t, tc.normalized, normalized)
			assert.Equal(t, tc.route, route)
		})
	}
}
