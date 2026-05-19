package storage

type Storage struct {
	Shares             []Share            `tree:"shares"`
	Queues             []Queue            `tree:"queues"`
	Accounts           []Account          `tree:"accounts"`
	ManagementPolicies []ManagementPolicy `tree:"management_policies"`
}

func (s *Storage) PostProcess() {
	// Link queues to accounts by StorageAccountName
	for i, queue := range s.Queues {
		for j := range s.Accounts {
			if queue.StorageAccountName.Value() == s.Accounts[j].Name.Value() {
				s.Queues[i].AccountReplicationType = s.Accounts[j].AccountReplicationType
				s.Queues[i].AccountKind = s.Accounts[j].AccountKind
				s.Queues[i].AccountTier = s.Accounts[j].AccountTier
				break
			}
		}
	}

	// Link shares to accounts by StorageAccountName
	for i, share := range s.Shares {
		for j := range s.Accounts {
			if share.StorageAccountName.Value() == s.Accounts[j].Name.Value() {
				s.Shares[i].AccountReplicationType = s.Accounts[j].AccountReplicationType
				s.Shares[i].AccountKind = s.Accounts[j].AccountKind
				s.Shares[i].AccountTier = s.Accounts[j].AccountTier
				break
			}
		}
	}

	// Link management policies to accounts by StorageAccountID
	for i, policy := range s.ManagementPolicies {
		for j := range s.Accounts {
			if policy.StorageAccountID.Value() == s.Accounts[j].ID {
				s.Accounts[j].Relationships.ManagementPolicy = &s.ManagementPolicies[i]
				break
			}
		}
	}
}
