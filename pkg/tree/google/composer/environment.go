package composer

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Environment struct {
	resource.Resource  `tree:"-"`
	ImageVersion       value.String `tree:"image_version"`
	EnvironmentSize    value.String `tree:"environment_size"`
	SchedulerCPU       value.Double `tree:"scheduler_cpu"`
	SchedulerMemoryGB  value.Double `tree:"scheduler_memory_gb"`
	SchedulerStorageGB value.Double `tree:"scheduler_storage_gb"`
	SchedulerCount     value.Int    `tree:"scheduler_count"`
	WebServerCPU       value.Double `tree:"web_server_cpu"`
	WebServerMemoryGB  value.Double `tree:"web_server_memory_gb"`
	WebServerStorageGB value.Double `tree:"web_server_storage_gb"`
	WorkerCPU          value.Double `tree:"worker_cpu"`
	WorkerMemoryGB     value.Double `tree:"worker_memory_gb"`
	WorkerStorageGB    value.Double `tree:"worker_storage_gb"`
	WorkerMinCount     value.Int    `tree:"worker_min_count"`
}
