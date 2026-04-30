package container

type Container struct {
	Registries                 []Registry                 `tree:"registries"`
	KubernetesClusters         []KubernetesCluster        `tree:"kubernetes_clusters"`
	KubernetesClusterNodePools []KubernetesClusterNodePool `tree:"kubernetes_cluster_node_pools"`
	AppEnvironments            []AppEnvironment           `tree:"app_environments"`
	Apps                       []App                      `tree:"apps"`
}

func (s *Container) PostProcess() {
	for i, nodePool := range s.KubernetesClusterNodePools {
		for j := range s.KubernetesClusters {
			if nodePool.ClusterID.Value() == s.KubernetesClusters[j].ID {
				s.KubernetesClusterNodePools[i].Relationships.Cluster = &s.KubernetesClusters[j]
				break
			}
		}
	}

	for i, app := range s.Apps {
		for j := range s.AppEnvironments {
			if app.ContainerAppEnvironmentID.Value() == s.AppEnvironments[j].ID {
				s.Apps[i].Relationships.Environment = &s.AppEnvironments[j]

				// link app to its workload profile within the environment
				for k := range s.AppEnvironments[j].WorkloadProfiles {
					if app.WorkloadProfileName.Value() == s.AppEnvironments[j].WorkloadProfiles[k].Name.Value() {
						s.Apps[i].Relationships.WorkloadProfile = &s.AppEnvironments[j].WorkloadProfiles[k]
						break
					}
				}
				break
			}
		}
	}
}
