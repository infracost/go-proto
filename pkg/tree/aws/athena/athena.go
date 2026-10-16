package athena

type Athena struct {
	Workgroups           []Workgroup           `tree:"workgroups"`
	CapacityReservations []CapacityReservation `tree:"capacity_reservations"`
}
