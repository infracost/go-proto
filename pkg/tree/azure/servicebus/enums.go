package servicebus

type NamespaceSKU uint32

const (
	NamespaceSKUUnknown NamespaceSKU = iota
	NamespaceSKUBasic
	NamespaceSKUStandard
	NamespaceSKUPremium
)
