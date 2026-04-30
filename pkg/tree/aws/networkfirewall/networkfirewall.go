package networkfirewall

type NetworkFirewall struct {
	Firewalls []Firewall `tree:"firewalls"`
}
