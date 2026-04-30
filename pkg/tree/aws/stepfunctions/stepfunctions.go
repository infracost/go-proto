package stepfunctions

type StepFunctions struct {
	StateMachines []StateMachine `tree:"state_machines"`
}
