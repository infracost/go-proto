package batch

type Batch struct {
	ComputeEnvironments []ComputeEnvironment `tree:"compute_environment"`
	JobDefinitions      []JobDefinition      `tree:"job_definition"`
}
