package docdb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ClusterInstance struct {
	resource.Resource `tree:"-"`
	InstanceClass     value.String `tree:"instance_class"`
}
