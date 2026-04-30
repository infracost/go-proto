package bigquery

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CapacityCommitment struct {
	resource.Resource `tree:"-"`
	SlotCount         value.Int                    `tree:"slot_count"`
	Plan              value.Value[CommitmentPlan]  `tree:"plan"`
}

type CommitmentPlan uint32

const (
	CommitmentPlanUnknown   CommitmentPlan = iota
	CommitmentPlanFlex
	CommitmentPlanMonthly
	CommitmentPlanAnnual
	CommitmentPlanThreeYear
)
