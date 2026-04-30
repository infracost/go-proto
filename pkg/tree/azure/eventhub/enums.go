package eventhub

type NamespaceSKU uint32

const (
	NamespaceSKUUnknown NamespaceSKU = iota
	NamespaceSKUBasic
	NamespaceSKUStandard
	NamespaceSKUPremium
)
