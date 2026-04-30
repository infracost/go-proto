package iothub

type IoTHub struct {
	Hubs []Hub `tree:"hubs"`
	DPSs []DPS `tree:"dpss"`
}
