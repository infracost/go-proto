package sesv2

type SESV2 struct {
	DedicatedIPAssignments []DedicatedIPAssignment `tree:"dedicated_ip_assignments"`
	DedicatedIPPools       []DedicatedIPPool       `tree:"dedicated_ip_pools"`
	EmailIdentities        []EmailIdentity         `tree:"email_identities"`
	ReceiptRules           []ReceiptRule           `tree:"receipt_rules"`
}
