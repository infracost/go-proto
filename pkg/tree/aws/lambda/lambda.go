package lambda

type Lambda struct {
	Functions                     []Function                     `tree:"functions"`
	ProvisionedConcurrencyConfigs []ProvisionedConcurrencyConfig `tree:"provisioned_concurrency_configs"`
	EventSourceMappings           []EventSourceMapping           `tree:"event_source_mappings"`
}
