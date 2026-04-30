package pubsub

type PubSub struct {
	Subscriptions []Subscription `tree:"subscriptions"`
	Topics        []Topic        `tree:"topics"`
}
