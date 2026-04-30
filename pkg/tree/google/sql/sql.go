package sql

type SQL struct {
	DatabaseInstances []DatabaseInstance `tree:"database_instances"`
}
