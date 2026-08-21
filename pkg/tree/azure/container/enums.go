package container

type RegistrySKU uint32

const (
	RegistrySKUUnknown RegistrySKU = iota
	RegistrySKUBasic
	RegistrySKUStandard
	RegistrySKUPremium
	RegistrySKUClassic
)

type KubernetesSKUTier uint32

const (
	KubernetesSKUTierUnknown KubernetesSKUTier = iota
	KubernetesSKUTierFree
	KubernetesSKUTierStandard
	KubernetesSKUTierPremium
)

type OSType uint32

const (
	OSTypeUnknown OSType = iota
	OSTypeLinux
	OSTypeWindows
)

type OSDiskType uint32

const (
	OSDiskTypeUnknown OSDiskType = iota
	OSDiskTypeManaged
	OSDiskTypeEphemeral
)
