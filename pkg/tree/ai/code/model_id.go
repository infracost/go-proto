package code

import "strings"

// routeSeparators are the SDKs that put a provider prefix in front of a model
// ID, and the separator each uses: PydanticAI's and LangChain init_chat_model's
// "provider:model", LiteLLM's and instructor's "provider/model". Other SDKs
// take the model ID as is, and their "/" and ":" are part of it (Ollama's
// "mistral:7b", OpenRouter's "anthropic/claude-3.5-sonnet" through the OpenAI
// SDK). A model built with one of PydanticAI's model classes
// (OllamaModel("mistral:7b")) takes the ID as its provider names it; the call
// records the class, and its ID is never normalized.
var routeSeparators = map[string]byte{
	"pydantic-ai": ':',
	"langchain":   ':',
	"litellm":     '/',
	"instructor":  '/',
}

// knownRoutes are the provider prefixes those SDKs use, after the parser's
// provider names. Only these are stripped, because "/" and ":" are also part
// of real model IDs (Together's "meta-llama/Llama-3-70b-chat-hf", Bedrock's
// "…-v1:0").
var knownRoutes = map[string]bool{
	"anthropic":                true,
	"anthropic_bedrock":        true,
	"anyscale":                 true,
	"azure":                    true,
	"azure-responses":          true,
	"azure_ai":                 true,
	"azure_openai":             true,
	"bedrock":                  true,
	"bedrock-mantle":           true,
	"bedrock_converse":         true,
	"bedrock_mantle_anthropic": true,
	"bedrock_mantle_openai":    true,
	"cerebras":                 true,
	"cohere":                   true,
	"cohere_chat":              true,
	"databricks":               true,
	"deepseek":                 true,
	"fireworks":                true,
	"fireworks_ai":             true,
	"gemini":                   true,
	"generative-ai":            true,
	"google":                   true,
	"google-cloud":             true,
	"google-gla":               true,
	"google-vertex":            true,
	"google_anthropic_vertex":  true,
	"google_genai":             true,
	"google_vertexai":          true,
	"groq":                     true,
	"hosted_vllm":              true,
	"huggingface":              true,
	"ibm":                      true,
	"mistral":                  true,
	"mistralai":                true,
	"nvidia":                   true,
	"nvidia_nim":               true,
	"ollama":                   true,
	"ollama_chat":              true,
	"openai":                   true,
	"openai-chat":              true,
	"openai-codex":             true,
	"openai-responses":         true,
	"openrouter":               true,
	"perplexity":               true,
	"pydantic-gateway":         true,
	"replicate":                true,
	"sagemaker":                true,
	"together":                 true,
	"together_ai":              true,
	"vertex_ai":                true,
	"vertex_ai_beta":           true,
	"vertexai":                 true,
	"watsonx":                  true,
	"writer":                   true,
	"xai":                      true,
}

// NormalizeModelID strips a known provider prefix from a model ID as written,
// when the call's SDK uses one, and returns the rest as the normalized ID with
// the prefix as the route. Anything else, including a Bedrock inference
// profile's region prefix ("us.anthropic…"), is kept as written.
//
// It doesn't know whether the ID was passed to a framework's model class,
// whose IDs are never normalized; ModelCall.NormalizedModelID applies that
// rule too.
func NormalizeModelID(sdk, id string) (normalized, route string) {
	sep, ok := routeSeparators[sdk]
	if !ok {
		return id, ""
	}
	route, rest, ok := strings.Cut(id, string(sep))
	if !ok || rest == "" || !knownRoutes[route] {
		return id, ""
	}
	// A route's own sub-route, such as LiteLLM's bedrock/converse/<model>,
	// is part of the route, not the model ID.
	for _, sub := range subRoutes[sdk+" "+route] {
		if model, ok := strings.CutPrefix(rest, sub+string(sep)); ok && model != "" {
			return model, route
		}
	}
	return rest, route
}

// subRoutes are the second segments an SDK's route can have before the model
// ID, keyed by "<sdk> <route>": LiteLLM's Bedrock route picks its API with
// bedrock/converse/<model> or bedrock/invoke/<model>.
var subRoutes = map[string][]string{
	"litellm bedrock": {"converse", "invoke"},
}

// NormalizedModelID returns the call's model ID with any provider prefix its
// SDK put in front of it removed, and that prefix as the route, as
// NormalizeModelID does. A model built with a framework's model class (see
// ModelClass) takes the ID as its provider names it, so it is returned as
// written, with no route.
//
// It reads ModelID's value only: an unresolved model is "", so check ModelID
// for flag.Synthetic first.
func (c ModelCall) NormalizedModelID() (normalized, route string) {
	id := c.ModelID.Value()
	if c.ModelClass.Value() != "" {
		return id, ""
	}
	return NormalizeModelID(c.SDK.Value(), id)
}
