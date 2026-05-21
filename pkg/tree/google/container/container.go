package container

type Container struct {
	Clusters   []Cluster   `tree:"clusters"`
	NodePools  []NodePool  `tree:"node_pools"`
	Registries []Registry  `tree:"registries"`
}

func (c *Container) PostProcess() {
	// Reset relationships first to keep this idempotent — PostProcess can run
	// multiple times (parser builds the tree, processor re-hydrates from proto
	// and runs PostProcess again).
	for j := range c.Clusters {
		c.Clusters[j].Relationships.NodePools = nil
		c.Clusters[j].Relationships.DefaultNodePool = nil
	}
	for i := range c.NodePools {
		c.NodePools[i].Relationships.Cluster = nil
	}

	// Link NodePools to Clusters by ClusterID
	for i, np := range c.NodePools {
		if np.ClusterID.IsEmpty() {
			continue
		}
		for j := range c.Clusters {
			cluster := &c.Clusters[j]
			if np.ClusterID.Equal(cluster.ID) || np.ClusterID.Value() == cluster.Name.Value() {
				c.NodePools[i].Relationships.Cluster = cluster

				// Inherit Zones from Cluster if not set
				if np.Zones.IsEmpty() {
					c.NodePools[i].Zones = cluster.Zones
				}

				// Append to Cluster's NodePools
				cluster.Relationships.NodePools = append(cluster.Relationships.NodePools, &c.NodePools[i])

				// Inline pools named "default-pool" (either an explicit block in HCL
				// or the synthetic default the parser creates when no explicit
				// default is present) become the cluster's DefaultNodePool.
				if np.IsChild && np.Name.Value() == "default-pool" && cluster.Relationships.DefaultNodePool == nil {
					cluster.Relationships.DefaultNodePool = &c.NodePools[i]
				}
				break
			}
		}
	}
}
