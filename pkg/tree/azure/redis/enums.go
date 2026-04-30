package redis

type CacheFamily uint32

const (
	CacheFamilyUnknown CacheFamily = iota
	CacheFamilyC
	CacheFamilyP
)
