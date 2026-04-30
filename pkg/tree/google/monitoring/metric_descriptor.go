package monitoring

import "github.com/infracost/go-proto/pkg/tree/resource"

type MetricDescriptor struct {
	resource.Resource `tree:"-"`
}
