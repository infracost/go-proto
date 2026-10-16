package athena

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CapacityReservation struct {
	resource.Resource `tree:"-"`
	TargetDPUs        value.Int `tree:"target_dpus"`
}
