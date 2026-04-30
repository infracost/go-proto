package redis

type Redis struct {
	Instances []Instance `tree:"instances"`
}
