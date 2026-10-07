package code

// Channel is the route a client's requests take to reach the model.
type Channel uint32

const (
	ChannelUnknown Channel = iota
	// ChannelDirect calls the model vendor's own API.
	ChannelDirect
	// ChannelBedrock calls the model through AWS Bedrock.
	ChannelBedrock
	// ChannelVertex calls the model through Google Vertex AI.
	ChannelVertex
	// ChannelAzureOpenAI calls the model through Azure OpenAI.
	ChannelAzureOpenAI
	// ChannelProxy calls the model through a gateway or proxy, such as
	// LiteLLM or OpenRouter. The Client's Gateway names which one.
	ChannelProxy
)
