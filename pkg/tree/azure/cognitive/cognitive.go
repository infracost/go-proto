package cognitive

type Cognitive struct {
	Accounts    []Account    `tree:"accounts"`
	Deployments []Deployment `tree:"deployments"`
	Languages   []Language   `tree:"languages"`
	LUIS        []LUIS       `tree:"luis"`
	Speeches    []Speech     `tree:"speeches"`
}

func (s *Cognitive) PostProcess() {
	for i, deployment := range s.Deployments {
		for j := range s.Accounts {
			if deployment.CognitiveAccountID.Value() == s.Accounts[j].ID {
				s.Deployments[i].Relationships.Account = &s.Accounts[j]
				break
			}
		}
	}
}
