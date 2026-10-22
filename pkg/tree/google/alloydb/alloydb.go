package alloydb

type AlloyDB struct {
	Clusters  []Cluster  `tree:"clusters"`
	Instances []Instance `tree:"instances"`
}
