package guardduty

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Detector struct {
	resource.Resource `tree:"-"`
	Enable            value.Bool `tree:"enable"`
}
