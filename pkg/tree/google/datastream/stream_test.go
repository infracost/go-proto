package datastream_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/google"
	"github.com/infracost/go-proto/pkg/tree/google/datastream"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamProtoRoundTrip(t *testing.T) {
	original := &tree.Tree{
		Google: google.Google{
			Datastream: datastream.Datastream{
				Streams: []datastream.Stream{{Resource: resource.Resource{ID: "stream"}}},
			},
		},
	}

	encoded, err := original.ToProto()
	require.NoError(t, err)

	decoded, err := tree.FromProto(encoded)
	require.NoError(t, err)
	require.Len(t, decoded.Google.Datastream.Streams, 1)
	assert.Equal(t, "stream", decoded.Google.Datastream.Streams[0].ID)
}
