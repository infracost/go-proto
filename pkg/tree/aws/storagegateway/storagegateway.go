package storagegateway

type StorageGateway struct {
	Gateways []Gateway `tree:"gateways"`
}
