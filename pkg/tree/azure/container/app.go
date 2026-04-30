package container

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type App struct {
	resource.Resource         `tree:"-"`
	ContainerAppEnvironmentID value.String `tree:"container_app_environment_id"`
	WorkloadProfileName       value.String `tree:"workload_profile_name"`
	CPUCores                  value.Double `tree:"cpu_cores"`
	MemoryGiB                 value.Double `tree:"memory_gib"`
	MinReplicas               value.Int    `tree:"min_replicas"`
	MaxReplicas               value.Int    `tree:"max_replicas"`

	Relationships AppRelationships `tree:"-"`
}

type AppRelationships struct {
	Environment     *AppEnvironment
	WorkloadProfile *WorkloadProfile
}
