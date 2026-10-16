package backup

type Backup struct {
	Vaults                   []Vault                   `tree:"vaults"`
	LogicallyAirGappedVaults []LogicallyAirGappedVault `tree:"logically_air_gapped_vaults"`
	RestoreTestingPlans      []RestoreTestingPlan      `tree:"restore_testing_plans"`
	Frameworks               []Framework               `tree:"frameworks"`
}
