package efs

type EFS struct {
	FileSystems               []FileSystem               `tree:"file_systems"`
	ReplicationConfigurations []ReplicationConfiguration `tree:"replication_configurations"`
}
