package glacier

type Glacier struct {
	Vaults []Vault `tree:"vaults"`
}
