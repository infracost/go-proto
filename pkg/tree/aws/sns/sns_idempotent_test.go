package sns

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &SNS{
		Topics: []Topic{
			{Resource: resource.Resource{ID: "topic-1"}},
		},
		Subscriptions: []Subscription{
			{
				Resource: resource.Resource{ID: "sub-1"},
				TopicARN: value.New("topic-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	subs := append([]*Subscription(nil), s.Topics[0].Relationships.Subscriptions...)
	assert.Len(t, subs, 1)

	s.PostProcess()
	assert.Equal(t, subs, s.Topics[0].Relationships.Subscriptions)
}
