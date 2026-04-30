package kms

type KMS struct {
	Keys         []Key         `tree:"keys"`
	ExternalKeys []ExternalKey `tree:"external_keys"`
}
