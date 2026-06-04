package dynamodb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DAXCluster struct {
	resource.Resource             `tree:"-"`
	ClusterEndpointEncryptionType value.String `tree:"cluster_endpoint_encryption_type"`
	ServerSideEncryptionEnabled   value.Bool   `tree:"server_side_encryption_enabled"`
}
