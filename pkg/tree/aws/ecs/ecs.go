package ecs

import "strings"

type ECS struct {
	TaskDefinitions          []TaskDefinition          `tree:"task_definitions"`
	Services                 []ServiceResource         `tree:"services"`
	Clusters                 []Cluster                 `tree:"clusters"`
	ClusterCapacityProviders []ClusterCapacityProviders `tree:"cluster_capacity_providers"`
	TaskSets                 []TaskSet                 `tree:"task_sets"`
	CapacityProviders        []CapacityProvider        `tree:"capacity_providers"`
}

func (e *ECS) PostProcess() {
	// Build capacity provider type map
	cpTypes := make(map[string]CapacityProviderType)
	for _, cp := range e.CapacityProviders {
		if name := cp.Name.Value(); name != "" {
			cpTypes[name] = cp.ProviderType.Value()
		}
	}

	// Link cluster capacity providers to clusters
	for i, cluster := range e.Clusters {
		for j := range e.ClusterCapacityProviders {
			cp := &e.ClusterCapacityProviders[j]
			if cluster.Name.Value() == cp.ClusterName.Value() || cluster.ID == cp.ClusterName.Value() {
				e.Clusters[i].Relationships.CapacityProviders = append(
					e.Clusters[i].Relationships.CapacityProviders, cp,
				)
			}
		}
		e.Clusters[i].Relationships.KnownCapacityProviderTypes = cpTypes
	}

	// Link task definitions to task sets
	for i, taskSet := range e.TaskSets {
		for j, taskDef := range e.TaskDefinitions {
			if taskSet.TaskDefinitionReference.Value() == taskDef.ID ||
				(taskDef.Family.Value() != "" && strings.HasPrefix(taskSet.TaskDefinitionReference.Value(), taskDef.Family.Value())) {
				e.TaskSets[i].Relationships.TaskDefinition = &e.TaskDefinitions[j]
				break
			}
		}
	}

	// Link clusters and task definitions to services
	for i, svc := range e.Services {
		e.Services[i].Relationships.KnownCapacityProviderTypes = cpTypes

		for j := range e.Clusters {
			c := &e.Clusters[j]
			if c.ID == svc.ClusterReference.Value() || c.Name.Value() == svc.ClusterReference.Value() {
				e.Services[i].Relationships.Cluster = c
				break
			}
		}

		for j := range e.TaskSets {
			ts := &e.TaskSets[j]
			if ts.Service.Value() == svc.ID || (!svc.Name.IsEmpty() && ts.Service.Value() == svc.Name.Value()) {
				e.Services[i].Relationships.TaskDefinition = ts.Relationships.TaskDefinition
				break
			}
		}

		for j, taskDef := range e.TaskDefinitions {
			if svc.TaskDefinitionReference.Value() == taskDef.ID ||
				(taskDef.Family.Value() != "" && strings.HasPrefix(svc.TaskDefinitionReference.Value(), taskDef.Family.Value())) {
				e.Services[i].Relationships.TaskDefinition = &e.TaskDefinitions[j]
				break
			}
		}
	}
}
