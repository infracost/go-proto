package lightsail

type Lightsail struct {
	Instances         []Instance         `tree:"instances"`
	Databases         []Database         `tree:"databases"`
	Disks             []Disk             `tree:"disks"`
	LoadBalancers     []LoadBalancer     `tree:"load_balancers"`
	ContainerServices []ContainerService `tree:"container_services"`
	Distributions     []Distribution     `tree:"distributions"`
	Buckets           []Bucket           `tree:"buckets"`
}
