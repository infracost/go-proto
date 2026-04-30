package notificationhub

type NotificationHub struct {
	Namespaces []Namespace `tree:"namespaces"`
}
