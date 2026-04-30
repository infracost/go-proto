package cosmosdb

type ConsistencyLevel uint32

const (
	ConsistencyLevelUnknown ConsistencyLevel = iota
	ConsistencyLevelStrong
	ConsistencyLevelBoundedStaleness
	ConsistencyLevelSession
	ConsistencyLevelConsistentPrefix
	ConsistencyLevelEventual
)

type CosmosDBKind uint32

const (
	CosmosDBKindUnknown CosmosDBKind = iota
	CosmosDBKindGlobalDocumentDB
	CosmosDBKindMongoDB
	CosmosDBKindParse
)

type BackupType uint32

const (
	BackupTypeUnknown BackupType = iota
	BackupTypePeriodic
	BackupTypeContinuous
)

type BackupStorageRedundancy uint32

const (
	BackupStorageRedundancyUnknown BackupStorageRedundancy = iota
	BackupStorageRedundancyGeo
	BackupStorageRedundancyLocal
	BackupStorageRedundancyZone
)
