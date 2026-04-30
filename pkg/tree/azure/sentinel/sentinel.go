package sentinel

type Sentinel struct {
	DataConnectors []DataConnector `tree:"data_connectors"`
}
