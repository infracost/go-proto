package dms

type DMS struct {
	ReplicationInstances []ReplicationInstance `tree:"replication_instances"`
	Endpoints            []Endpoint            `tree:"endpoints"`
}
