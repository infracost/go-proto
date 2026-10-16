package redshiftserverless

type RedshiftServerless struct {
	Workgroups []Workgroup `tree:"workgroups"`
	Namespaces []Namespace `tree:"namespaces"`
	Snapshots  []Snapshot  `tree:"snapshots"`
}
