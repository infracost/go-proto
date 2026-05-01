package elasticache

type CacheEngine uint32

const (
	CacheEngineUnknown   CacheEngine = iota
	CacheEngineMemcached
	CacheEngineRedis
	CacheEngineValkey
)
