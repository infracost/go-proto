package bedrock

import "github.com/infracost/go-proto/pkg/tree/resource"

type AgentCoreRuntime struct {
	resource.Resource `tree:"-"`
}

type AgentCoreBrowser struct {
	resource.Resource `tree:"-"`
}

type AgentCoreCodeInterpreter struct {
	resource.Resource `tree:"-"`
}

type AgentCoreGateway struct {
	resource.Resource `tree:"-"`
}

type AgentCoreMemory struct {
	resource.Resource `tree:"-"`
}
