package code

// Channel is the route a client's requests take to reach the model: direct to
// the model's vendor, or through a middleman.
//
// The values are pinned, not iota: 2, 3 and 4 were the retired Bedrock,
// Vertex and Azure OpenAI channels, which are now ChannelProxy with a Gateway,
// and must never be reused.
type Channel uint32

const (
	ChannelUnknown Channel = 0
	// ChannelDirect calls the model vendor's own API, or a self-hosted
	// inference server such as Ollama or vLLM that runs the model itself.
	ChannelDirect Channel = 1
	// ChannelProxy calls the model through a middleman: a cloud platform such
	// as Bedrock, Vertex AI or Azure OpenAI, or a gateway such as OpenRouter
	// or LiteLLM. The Client's Gateway names which one.
	ChannelProxy Channel = 5
)
