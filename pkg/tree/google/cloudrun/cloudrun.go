package cloudrun

type CloudRun struct {
	V1Services []V1Service `tree:"v1_services"`
	V2Jobs     []V2Job     `tree:"v2_jobs"`
	V2Services []V2Service `tree:"v2_services"`
}
