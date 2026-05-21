package logging

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

// Sink represents google_logging_*_sink resources. They share the same
// "Log Volume" pricing as bucket configs — the resource itself is free and the
// only cost is the data volume measured by usage.
type Sink struct {
	resource.Resource `tree:"-"`
}
