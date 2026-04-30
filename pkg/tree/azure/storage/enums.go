package storage

type AccessTier uint32

const (
	AccessTierUnknown AccessTier = iota
	AccessTierHot
	AccessTierCool
	AccessTierCold
	AccessTierPremium
)

type AccountKind uint32

const (
	AccountKindUnknown AccountKind = iota
	AccountKindBlobStorage
	AccountKindBlockBlobStorage
	AccountKindFileStorage
	AccountKindStorage
	AccountKindStorageV2
)

type AccountReplicationType uint32

const (
	AccountReplicationTypeUnknown AccountReplicationType = iota
	AccountReplicationTypeLRS
	AccountReplicationTypeGRS
	AccountReplicationTypeRAGRS
	AccountReplicationTypeZRS
	AccountReplicationTypeGZRS
	AccountReplicationTypeRAGZRS
)

type AccountTier uint32

const (
	AccountTierUnknown AccountTier = iota
	AccountTierStandard
	AccountTierPremium
)
