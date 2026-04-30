package loganalytics

type LogAnalytics struct {
	Workspaces []Workspace `tree:"workspaces"`
}
