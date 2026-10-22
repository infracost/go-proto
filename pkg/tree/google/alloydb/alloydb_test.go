package alloydb_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/google"
	"github.com/infracost/go-proto/pkg/tree/google/alloydb"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlloyDBProtoRoundTrip(t *testing.T) {
	original := &tree.Tree{
		Google: google.Google{
			AlloyDB: alloydb.AlloyDB{
				Clusters: []alloydb.Cluster{{Resource: resource.Resource{ID: "cluster"}}},
				Instances: []alloydb.Instance{{
					InstanceType:     value.New("READ_POOL", 0, "instance_type", nil),
					AvailabilityType: value.New("ZONAL", 0, "availability_type", nil),
					CPUCount:         value.New(int64(4), 0, "machine_config.0.cpu_count", nil),
					NodeCount:        value.New(int64(3), 0, "read_pool_config.0.node_count", nil),
				}},
			},
		},
	}

	encoded, err := original.ToProto()
	require.NoError(t, err)

	decoded, err := tree.FromProto(encoded)
	require.NoError(t, err)
	require.Len(t, decoded.Google.AlloyDB.Clusters, 1)
	require.Len(t, decoded.Google.AlloyDB.Instances, 1)

	instance := decoded.Google.AlloyDB.Instances[0]
	assert.Equal(t, "READ_POOL", instance.InstanceType.Value())
	assert.Equal(t, "ZONAL", instance.AvailabilityType.Value())
	assert.Equal(t, int64(4), instance.CPUCount.Value())
	assert.Equal(t, int64(3), instance.NodeCount.Value())
}
