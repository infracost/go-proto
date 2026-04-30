package secretmanager

type SecretManager struct {
	Secrets        []Secret        `tree:"secrets"`
	SecretVersions []SecretVersion `tree:"secret_versions"`
}

func (s *SecretManager) PostProcess() {
	for i, sv := range s.SecretVersions {
		for j, secret := range s.Secrets {
			if sv.SecretRef.Value() == secret.ID || sv.SecretRef.Value() == secret.Name.Value() {
				s.SecretVersions[i].Relationships.Secret = &s.Secrets[j]
				if sv.ReplicationLocations.IsDefaultOrEmpty() {
					s.SecretVersions[i].ReplicationLocations = secret.ReplicationLocations
				}
				break
			}
		}
	}
}
