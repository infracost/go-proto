package synapse

type SparkPoolNodeSize uint32

const (
	SparkPoolNodeSizeUnknown SparkPoolNodeSize = iota
	SparkPoolNodeSizeSmall
	SparkPoolNodeSizeMedium
	SparkPoolNodeSizeLarge
	SparkPoolNodeSizeXLarge
	SparkPoolNodeSizeXXLarge
	SparkPoolNodeSizeXXXLarge
)
