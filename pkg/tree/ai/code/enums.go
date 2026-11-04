package code

// Channel is the route a client's requests take to reach the model: direct to
// the model's vendor, or through a middleman.
type Channel uint32

const (
	ChannelUnknown Channel = iota
	// ChannelDirect calls the model vendor's own API, or a self-hosted
	// inference server such as Ollama or vLLM that runs the model itself.
	ChannelDirect
	// ChannelProxy calls the model through a middleman: a cloud platform such
	// as Bedrock, Vertex AI or Azure OpenAI, or a gateway such as OpenRouter
	// or LiteLLM. The Client's Gateway names which one.
	ChannelProxy
)
