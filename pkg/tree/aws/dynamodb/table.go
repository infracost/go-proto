package dynamodb

import (
	"github.com/infracost/go-proto/pkg/tree/aws/appautoscaling"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Table struct {
	resource.Resource          `tree:"-"`
	Name                       value.String             `tree:"name"`
	BillingMode                value.Value[BillingMode] `tree:"billing_mode"`
	WriteCapacity              value.Int                `tree:"write_capacity"`
	ReadCapacity               value.Int                `tree:"read_capacity"`
	PointInTimeRecoveryEnabled value.Bool               `tree:"point_in_time_recovery_enabled"`
	ReplicaRegions             value.List[string]       `tree:"replica_regions"`
	PropagateTags              value.Bool               `tree:"propagate_tags"`
	TimeToLiveSpecs            []TTLSpec                `tree:"time_to_live_specs"`

	Relationships TableRelationships `tree:"-"`
}

type TTLSpec struct {
	Attribute value.String `tree:"attribute"`
	Enabled   value.Bool   `tree:"enabled"`
}

type TableRelationships struct {
	AppAutoscalingTargets []appautoscaling.Target
}

type BillingMode uint32

const (
	BillingModeUnknown BillingMode = iota
	BillingModeProvisioned
	BillingModePayPerRequest
)
