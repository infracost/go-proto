package appconfiguration

type AppConfiguration struct {
	Configurations []Configuration `tree:"configurations"`
}
