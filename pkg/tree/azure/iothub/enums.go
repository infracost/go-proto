package iothub

type HubSKU uint32

const (
	HubSKUUnknown HubSKU = iota
	HubSKUB1
	HubSKUB2
	HubSKUB3
	HubSKUF1
	HubSKUS1
	HubSKUS2
	HubSKUS3
	HubSKUGEN2
)

type DPSSKU uint32

const (
	DPSSKUUnknown DPSSKU = iota
	DPSSKUS1
)
