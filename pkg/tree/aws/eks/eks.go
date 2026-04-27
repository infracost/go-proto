package eks

type EKS struct {
	NodeGroups      []NodeGroup      `tree:"node_groups"`
	Clusters        []Cluster        `tree:"clusters"`
	FargateProfiles []FargateProfile `tree:"fargate_profiles"`
}
