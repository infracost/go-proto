package monitoring

type Monitoring struct {
	MetricDescriptors []MetricDescriptor `tree:"metric_descriptors"`
}
