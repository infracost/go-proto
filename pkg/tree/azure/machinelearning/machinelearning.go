package machinelearning

type MachineLearning struct {
	ComputeClusters   []ComputeCluster  `tree:"compute_clusters"`
	ComputeInstances  []ComputeInstance `tree:"compute_instances"`
}
