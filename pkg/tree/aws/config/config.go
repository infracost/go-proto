package config

type Config struct {
	ConfigRules            []ConfigRule            `tree:"config_rules"`
	ConfigurationRecorders []ConfigurationRecorder `tree:"configuration_recorders"`
	ConformancePacks       []ConformancePack       `tree:"conformance_packs"`
}
