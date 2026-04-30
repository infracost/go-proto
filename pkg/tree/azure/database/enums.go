package database

type SQLLicenseType uint32

const (
	SQLLicenseTypeUnknown SQLLicenseType = iota
	SQLLicenseTypeLicenseIncluded
	SQLLicenseTypeBasePrice
)

type BackupStorageType uint32

const (
	BackupStorageTypeUnknown BackupStorageType = iota
	BackupStorageTypeGeo
	BackupStorageTypeLocal
	BackupStorageTypeZone
)

type StorageAccountType uint32

const (
	StorageAccountTypeUnknown StorageAccountType = iota
	StorageAccountTypeGRS
	StorageAccountTypeLRS
	StorageAccountTypeZRS
)
