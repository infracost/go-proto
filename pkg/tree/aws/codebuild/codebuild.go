package codebuild

type CodeBuild struct {
	Projects []Project `tree:"projects"`
	Fleets   []Fleet   `tree:"fleets"`
}
