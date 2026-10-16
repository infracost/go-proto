package timestream

type Timestream struct {
	InfluxDBInstances []InfluxDBInstance `tree:"influxdb_instances"`
	InfluxDBClusters  []InfluxDBCluster  `tree:"influxdb_clusters"`
	Tables            []Table            `tree:"tables"`
}
