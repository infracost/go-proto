package firewall

type Firewall struct {
	Policies         []Policy         `tree:"policies"`
	CollectionGroups []CollectionGroup `tree:"collection_groups"`
}

func (s *Firewall) PostProcess() {
	for i, group := range s.CollectionGroups {
		for j := range s.Policies {
			if group.FirewallPolicyID.Value() == s.Policies[j].ID {
				s.Policies[j].Relationships.CollectionGroups = append(s.Policies[j].Relationships.CollectionGroups, s.CollectionGroups[i])
				break
			}
		}
	}
}
