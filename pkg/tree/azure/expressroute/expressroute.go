package expressroute

type ExpressRoute struct {
	Connections []Connection `tree:"connections"`
	Gateways    []Gateway    `tree:"gateways"`
}
