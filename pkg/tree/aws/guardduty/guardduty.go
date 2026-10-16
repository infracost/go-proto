package guardduty

type GuardDuty struct {
	Detectors              []Detector              `tree:"detectors"`
	MalwareProtectionPlans []MalwareProtectionPlan `tree:"malware_protection_plans"`
}
