package redshift

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource `tree:"-"`
	NodeType          value.String `tree:"node_type"`
	NumberOfNodes     value.Int    `tree:"number_of_nodes"`
	Encrypted         value.Bool   `tree:"encrypted"`
}
