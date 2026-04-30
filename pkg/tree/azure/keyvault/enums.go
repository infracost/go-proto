package keyvault

type VaultSKU uint32

const (
	VaultSKUUnknown VaultSKU = iota
	VaultSKUStandard
	VaultSKUPremium
)

type KeyType uint32

const (
	KeyTypeUnknown KeyType = iota
	KeyTypeEC
	KeyTypeECHSM
	KeyTypeRSA
	KeyTypeRSAHSM
	KeyTypeOct
)
