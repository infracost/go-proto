package cloudformation

type CloudFormation struct {
	Stacks    []Stack    `tree:"stacks"`
	StackSets []StackSet `tree:"stack_sets"`
}
