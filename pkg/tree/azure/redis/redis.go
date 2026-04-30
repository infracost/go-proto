package redis

type Redis struct {
	Caches []Cache `tree:"caches"`
}
