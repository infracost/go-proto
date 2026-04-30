package sagemaker

type SageMaker struct {
	NotebookInstances []NotebookInstance `tree:"notebook_instances"`
}
