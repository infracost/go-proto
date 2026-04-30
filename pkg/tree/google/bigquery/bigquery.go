package bigquery

type BigQuery struct {
	Datasets           []Dataset           `tree:"datasets"`
	Tables             []Table             `tree:"tables"`
	Reservations       []Reservation       `tree:"reservations"`
	CapacityCommitments []CapacityCommitment `tree:"capacity_commitments"`
}
