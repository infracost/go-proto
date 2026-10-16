package route53recoverycontrolconfig

import "github.com/infracost/go-proto/pkg/tree/resource"

// ReadinessCheck is a Route 53 ARC readiness check (route53recoveryreadiness).
type ReadinessCheck struct {
	resource.Resource `tree:"-"`
}
