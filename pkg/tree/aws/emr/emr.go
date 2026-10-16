package emr

type EMR struct {
	Clusters                  []Cluster                  `tree:"cluster"`
	InstanceGroups            []InstanceGroup            `tree:"instance_group"`
	InstanceFleets            []InstanceFleet            `tree:"instance_fleet"`
	ServerlessApplications    []ServerlessApplication    `tree:"serverless_applications"`
	ContainersVirtualClusters []ContainersVirtualCluster `tree:"containers_virtual_clusters"`
}
