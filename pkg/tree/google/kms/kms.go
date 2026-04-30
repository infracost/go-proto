package kms

type KMS struct {
	CryptoKeys []CryptoKey `tree:"crypto_keys"`
}
