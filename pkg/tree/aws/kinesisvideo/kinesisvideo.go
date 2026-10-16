package kinesisvideo

type KinesisVideo struct {
	Streams           []Stream           `tree:"streams"`
	SignalingChannels []SignalingChannel `tree:"signaling_channels"`
}
