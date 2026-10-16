package route53recoverycontrolconfig

type Route53RecoveryControlConfig struct {
	Clusters        []Cluster        `tree:"clusters"`
	ReadinessChecks []ReadinessCheck `tree:"readiness_checks"`
}
