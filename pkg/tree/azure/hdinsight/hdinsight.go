package hdinsight

type HDInsight struct {
	HadoopClusters           []HadoopCluster           `tree:"hadoop_clusters"`
	HBaseClusters            []HBaseCluster             `tree:"hbase_clusters"`
	SparkClusters            []SparkCluster             `tree:"spark_clusters"`
	InteractiveQueryClusters []InteractiveQueryCluster  `tree:"interactive_query_clusters"`
	KafkaClusters            []KafkaCluster             `tree:"kafka_clusters"`
}
