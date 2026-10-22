package compute_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/google"
	"github.com/infracost/go-proto/pkg/tree/google/compute"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataflowFlexTemplateJobProtoRoundTrip(t *testing.T) {
	original := &tree.Tree{
		Google: google.Google{
			Compute: compute.Compute{
				DataflowFlexTemplateJobs: []compute.DataflowFlexTemplateJob{{
					MachineType: value.New("n1-standard-4", 0, "machine_type", nil),
					NumWorkers:  value.New(int64(3), 0, "num_workers", nil),
				}},
			},
		},
	}

	encoded, err := original.ToProto()
	require.NoError(t, err)

	decoded, err := tree.FromProto(encoded)
	require.NoError(t, err)
	require.Len(t, decoded.Google.Compute.DataflowFlexTemplateJobs, 1)

	job := decoded.Google.Compute.DataflowFlexTemplateJobs[0]
	assert.Equal(t, "n1-standard-4", job.MachineType.Value())
	assert.Equal(t, int64(3), job.NumWorkers.Value())
}
