package lightsail

type Lightsail struct {
	Instances []Instance `tree:"instances"`
}
