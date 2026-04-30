package sns

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Topic struct {
	resource.Resource `tree:"-"`
	FIFO              value.Bool `tree:"fifo"`

	Relationships TopicRelationships `tree:"-"`
}

type TopicRelationships struct {
	Subscriptions []*Subscription
}
