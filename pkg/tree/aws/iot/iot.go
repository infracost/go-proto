package iot

type IoT struct {
	Things     []Thing     `tree:"things"`
	TopicRules []TopicRule `tree:"topic_rules"`
}
