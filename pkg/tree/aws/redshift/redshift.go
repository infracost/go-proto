package redshift

type Redshift struct {
	Clusters []Cluster `tree:"clusters"`
}
