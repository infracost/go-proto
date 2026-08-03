package recoveryservices

type StorageModeType uint32

const (
	StorageModeTypeUnknown StorageModeType = iota
	StorageModeTypeGeoRedundant
	StorageModeTypeLocallyRedundant
	StorageModeTypeZoneRedundant
	StorageModeTypeReadAccessGeoZoneRedundant
)
