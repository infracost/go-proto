package integration

type Integration struct {
	ServiceEnvironments []ServiceEnvironment `tree:"service_environments"`
}
