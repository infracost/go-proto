package batch

type Batch struct {
	ComputeEnvironemnts []ComputeEnfironment `tree:"comupte_environment"`
	JobDefinitions      []JobDefinition      `tree:"job_definition"`
}
