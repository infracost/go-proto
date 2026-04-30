package servicenetworking

type ServiceNetworking struct {
	Connections []Connection `tree:"connections"`
}
