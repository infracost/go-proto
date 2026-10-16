package msk

type MSK struct {
	Clusters           []Cluster           `tree:"clusters"`
	ServerlessClusters []ServerlessCluster `tree:"serverless_clusters"`
	Connectors         []Connector         `tree:"connectors"`
	Replicators        []Replicator        `tree:"replicators"`
	VPCConnections     []VPCConnection     `tree:"vpc_connections"`
}
