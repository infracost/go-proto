package redis

type Redis struct {
	Instances            []Instance            `tree:"instances"`
	MemorystoreInstances []MemorystoreInstance `tree:"memorystore_instances"`
}
