package ssm

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Activation struct {
	resource.Resource `tree:"-"`
	RegistrationLimit value.Int `tree:"registration_limit"`
}
