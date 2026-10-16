package grafana

type Grafana struct {
	Workspaces          []Workspace          `tree:"workspaces"`
	LicenseAssociations []LicenseAssociation `tree:"license_associations"`
}

func (g *Grafana) PostProcess() {
	for i := range g.Workspaces {
		g.Workspaces[i].Relationships = WorkspaceRelationships{}
	}
	for i := range g.LicenseAssociations {
		for j := range g.Workspaces {
			if g.Workspaces[j].ID == g.LicenseAssociations[i].WorkspaceID.Value() {
				g.Workspaces[j].Relationships.LicenseAssociation = &g.LicenseAssociations[i]
				break
			}
		}
	}
}
