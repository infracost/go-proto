package eventhub

type EventHub struct {
	Namespaces []Namespace `tree:"namespaces"`
}
