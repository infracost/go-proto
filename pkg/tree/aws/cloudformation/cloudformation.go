package cloudformation

type CloudFormation struct {
	Stacks    []Stack    `tree:"stacks,flatten"`
	StackSets []StackSet `tree:"stack_sets,flatten"`
}
