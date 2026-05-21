package frontdoor

type FrontDoor struct {
	FirewallPolicies []FirewallPolicy `tree:"firewall_policies"`
	Classics         []Classic        `tree:"classics"`
}
