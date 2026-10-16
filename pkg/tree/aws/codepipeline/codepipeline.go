package codepipeline

type CodePipeline struct {
	Pipelines []Pipeline `tree:"pipelines"`
}
