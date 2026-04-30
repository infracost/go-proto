package bigquery

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Reservation struct {
	resource.Resource `tree:"-"`
	SlotCapacity      value.Int `tree:"slot_capacity"`
}
