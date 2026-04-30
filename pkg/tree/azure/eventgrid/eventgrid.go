package eventgrid

type EventGrid struct {
	Topics []Topic `tree:"topics"`
}
