package kms

type KMS struct {
	Keys                []Key                `tree:"keys"`
	ExternalKeys        []ExternalKey        `tree:"external_keys"`
	ReplicaKeys         []ReplicaKey         `tree:"replica_keys"`
	ReplicaExternalKeys []ReplicaExternalKey `tree:"replica_external_keys"`
}
