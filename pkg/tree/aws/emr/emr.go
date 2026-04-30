package emr

type EMR struct {
	Clusters       []Cluster       `tree:"cluster"`
	InstanceGroups []InstanceGroup `tree:"instance_group"`
	InstanceFleets []InstanceFleet `tree:"instance_fleet"`
}
