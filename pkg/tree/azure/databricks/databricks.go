package databricks

type Databricks struct {
	Workspaces []Workspace `tree:"workspaces"`
}
