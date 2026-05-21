package logging

type Logging struct {
	BucketConfigs []BucketConfig `tree:"bucket_configs"`
	Sinks         []Sink         `tree:"sinks"`
}
