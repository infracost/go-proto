package docdb

type DocDB struct {
	Clusters         []Cluster         `tree:"clusters"`
	ClusterInstances []ClusterInstance `tree:"cluster_instances"`
	ClusterSnapshots []ClusterSnapshot `tree:"cluster_snapshots"`
}
