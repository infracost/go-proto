package medialive

type MediaLive struct {
	Channels []Channel `tree:"channels"`
}
