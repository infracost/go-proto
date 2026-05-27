package recoveryservices

type RecoveryServices struct {
	Vaults             []Vault             `tree:"vaults"`
	BackupProtectedVMs []BackupProtectedVM `tree:"backup_protected_vms"`
}

func (s *RecoveryServices) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range s.Vaults {
		s.Vaults[i].Relationships.ProtectedVMs = nil
	}

	for i, vm := range s.BackupProtectedVMs {
		for j := range s.Vaults {
			if vm.RecoveryVaultName.Value() == s.Vaults[j].Name.Value() {
				s.Vaults[j].Relationships.ProtectedVMs = append(s.Vaults[j].Relationships.ProtectedVMs, &s.BackupProtectedVMs[i])
				break
			}
		}
	}
}
