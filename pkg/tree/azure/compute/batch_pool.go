package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BatchPool struct {
	resource.Resource `tree:"-"`
	VMSize            value.String `tree:"vm_size"`
	NumberOfNodes     value.Int    `tree:"number_of_nodes"`
}
