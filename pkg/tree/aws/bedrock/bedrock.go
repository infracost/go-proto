package bedrock

type Bedrock struct {
	ProvisionedModelThroughputs []ProvisionedModelThroughput `tree:"provisioned_model_throughputs"`
	CustomModels                []CustomModel                `tree:"custom_models"`
	Guardrails                  []Guardrail                  `tree:"guardrails"`
	Flows                       []Flow                       `tree:"flows"`
	AgentCoreRuntimes           []AgentCoreRuntime           `tree:"agent_core_runtimes"`
	AgentCoreBrowsers           []AgentCoreBrowser           `tree:"agent_core_browsers"`
	AgentCoreCodeInterpreters   []AgentCoreCodeInterpreter   `tree:"agent_core_code_interpreters"`
	AgentCoreGateways           []AgentCoreGateway           `tree:"agent_core_gateways"`
	AgentCoreMemories           []AgentCoreMemory            `tree:"agent_core_memories"`
}
