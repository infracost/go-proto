package dynamodb

type DynamoDB struct {
	Tables                       []Table                       `tree:"tables"`
	DAXClusters                  []DAXCluster                  `tree:"dax_clusters"`
	TableReplicas                []TableReplica                `tree:"table_replicas"`
	KinesisStreamingDestinations []KinesisStreamingDestination `tree:"kinesis_streaming_destinations"`
	ContributorInsights          []ContributorInsights         `tree:"contributor_insights"`
}
