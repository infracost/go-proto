package memorydb

type MemoryDB struct {
	Clusters []Cluster `tree:"clusters"`
}
