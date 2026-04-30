package artifactregistry

type ArtifactRegistry struct {
	Repositories []Repository `tree:"repositories"`
}
