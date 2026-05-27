package kinesis

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	k := &Kinesis{
		Streams: []Stream{
			{Resource: resource.Resource{ID: "stream-1"}},
		},
		DeliveryStreams: []FirehoseDeliveryStream{
			{
				Resource:         resource.Resource{ID: "ds-1"},
				KinesisStreamARN: value.New("arn:aws:kinesis:us-east-1:123:stream/stream-1", 0, "", nil),
			},
		},
	}

	k.PostProcess()
	stream := k.DeliveryStreams[0].Relationships.KinesisStream

	k.PostProcess()
	assert.Equal(t, stream, k.DeliveryStreams[0].Relationships.KinesisStream)
	assert.NotNil(t, stream)
}
