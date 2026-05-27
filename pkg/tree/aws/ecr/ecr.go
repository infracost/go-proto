package ecr

type ECR struct {
	Repositories      []Repository      `tree:"repositories"`
	LifecyclePolicies []LifecyclePolicy `tree:"lifecycle_policies"`
}

func (e *ECR) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range e.Repositories {
		e.Repositories[i].Relationships = RepositoryRelationships{}
	}

	for i, repo := range e.Repositories {
		for j, policy := range e.LifecyclePolicies {
			if !policy.RepositoryName.IsEmpty() && policy.RepositoryName.Value() == repo.Name.Value() {
				e.Repositories[i].Relationships.LifecyclePolicies = append(
					e.Repositories[i].Relationships.LifecyclePolicies,
					&e.LifecyclePolicies[j],
				)
			}
		}
	}
}
