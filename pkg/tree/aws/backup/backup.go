package backup

type Backup struct {
	Vaults []Vault `tree:"vaults"`
}
