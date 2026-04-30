package keyvault

type KeyVault struct {
	Vaults       []Vault       `tree:"vaults"`
	Keys         []Key         `tree:"keys"`
	ManagedHSMs  []ManagedHSM  `tree:"managed_hsms"`
	Certificates []Certificate `tree:"certificates"`
}

func (s *KeyVault) PostProcess() {
	for i, key := range s.Keys {
		for j := range s.Vaults {
			if key.KeyVaultID.Value() == s.Vaults[j].ID {
				s.Keys[i].SKUName = s.Vaults[j].SKUName
				break
			}
		}
	}
}
