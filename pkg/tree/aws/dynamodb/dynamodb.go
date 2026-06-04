package dynamodb

type DynamoDB struct {
	Tables      []Table      `tree:"tables"`
	DAXClusters []DAXCluster `tree:"dax_clusters"`
}
