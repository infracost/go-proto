package recoveryservices

type RecoveryServices struct {
	Vaults             []Vault             `tree:"vaults"`
	BackupProtectedVMs []BackupProtectedVM `tree:"backup_protected_vms"`
}

func (s *RecoveryServices) PostProcess() {
	for i, vm := range s.BackupProtectedVMs {
		for j := range s.Vaults {
			if vm.RecoveryVaultName.Value() == s.Vaults[j].Name.Value() {
				s.Vaults[j].Relationships.ProtectedVMs = append(s.Vaults[j].Relationships.ProtectedVMs, &s.BackupProtectedVMs[i])
				break
			}
		}
	}
}
