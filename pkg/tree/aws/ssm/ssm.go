package ssm

type SSM struct {
	Activations []Activation `tree:"activations"`
	Parameters  []Parameter  `tree:"parameters"`
}
