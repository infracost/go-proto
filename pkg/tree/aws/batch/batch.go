package batch

type Batch struct {
	ComputeEnvironemnts []ComputeEnvironment `tree:"compute_environment"`
	JobDefinitions      []JobDefinition      `tree:"job_definition"`
}
