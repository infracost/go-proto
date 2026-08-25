package redis

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MemorystoreInstance struct {
	resource.Resource      `tree:"-"`
	NodeType               value.Value[MemorystoreNodeType] `tree:"node_type"`
	ShardCount             value.Int                        `tree:"shard_count"`
	ReplicaCount           value.Int                        `tree:"replica_count"`
	PersistenceMode        value.String                     `tree:"persistence_mode"`
	AutomatedBackupEnabled value.Bool                       `tree:"automated_backup_enabled"`
}

type MemorystoreNodeType uint32

const (
	MemorystoreNodeTypeUnknown MemorystoreNodeType = iota
	MemorystoreNodeTypeSharedCoreNano
	MemorystoreNodeTypeCustomPico
	MemorystoreNodeTypeCustomMicro
	MemorystoreNodeTypeCustomMini
	MemorystoreNodeTypeStandardSmall
	MemorystoreNodeTypeHighmemMedium
	MemorystoreNodeTypeHighcpuMedium
	MemorystoreNodeTypeStandardLarge
	MemorystoreNodeTypeHighmemXLarge
	MemorystoreNodeTypeHighmem2XLarge
)
