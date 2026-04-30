package compute

type DiskStorageType uint32

const (
	DiskStorageTypeUnknown      DiskStorageType = iota
	DiskStorageTypeStandardLRS
	DiskStorageTypeStandardSSDLRS
	DiskStorageTypeStandardSSDZRS
	DiskStorageTypePremiumLRS
	DiskStorageTypePremiumZRS
	DiskStorageTypePremiumV2LRS
	DiskStorageTypeUltraSSDLRS
)

type OSType uint32

const (
	OSTypeUnknown OSType = iota
	OSTypeLinux
	OSTypeWindows
)

type LicenseType uint32

const (
	LicenseTypeUnknown      LicenseType = iota
	LicenseTypeNone
	LicenseTypeWindowsClient
	LicenseTypeWindowsServer
	LicenseTypeRHELBYOS
	LicenseTypeSLESBYOS
)
