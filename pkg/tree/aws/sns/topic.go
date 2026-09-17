package sns

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Topic struct {
	resource.Resource `tree:"-"`
	FIFO              value.Bool `tree:"fifo"`
	// ArchiveRetentionDays is 0 when message archiving is disabled.
	ArchiveRetentionDays value.Int `tree:"archive_retention_days"`

	Relationships TopicRelationships `tree:"-"`
}

type TopicRelationships struct {
	Subscriptions []*Subscription
}
