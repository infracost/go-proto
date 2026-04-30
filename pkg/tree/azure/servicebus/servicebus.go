package servicebus

type ServiceBus struct {
	Namespaces []Namespace `tree:"namespaces"`
}
