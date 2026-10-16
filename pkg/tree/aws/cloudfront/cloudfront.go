package cloudfront

type CloudFront struct {
	Distributions      []Distribution      `tree:"distributions"`
	Functions          []Function          `tree:"functions"`
	RealtimeLogConfigs []RealtimeLogConfig `tree:"realtime_log_configs"`
	KeyValueStores     []KeyValueStore     `tree:"key_value_stores"`
}
