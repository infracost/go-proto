package datafactory

type ADFComputeType uint32

const (
	ADFComputeTypeUnknown ADFComputeType = iota
	ADFComputeTypeGeneral
	ADFComputeTypeComputeOptimized
	ADFComputeTypeMemoryOptimized
)

type SSISEdition uint32

const (
	SSISEditionUnknown SSISEdition = iota
	SSISEditionStandard
	SSISEditionEnterprise
)

type SSISLicenseType uint32

const (
	SSISLicenseTypeUnknown          SSISLicenseType = iota
	SSISLicenseTypeLicenseIncluded
	SSISLicenseTypeBasePrice
)
