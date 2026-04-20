package cloudformation

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type StackSet struct {
	resource.Resource `tree:"-"`
	TemplateBody      value.String `tree:"template_body"`
}
