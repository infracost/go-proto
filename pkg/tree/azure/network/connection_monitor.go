package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ConnectionMonitor struct {
	resource.Resource `tree:"-"`
	Tests             value.Int `tree:"tests"`
}
