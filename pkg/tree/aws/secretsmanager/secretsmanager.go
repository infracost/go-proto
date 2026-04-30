package secretsmanager

type SecretsManager struct {
	Secrets []Secret `tree:"secrets"`
}
