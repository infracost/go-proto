package dns

type DNS struct {
	ManagedZones []ManagedZone `tree:"managed_zones"`
	RecordSets   []RecordSet   `tree:"record_sets"`
}
