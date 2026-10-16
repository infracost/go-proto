package sagemaker

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type App struct {
	resource.Resource `tree:"-"`
	AppType           value.String `tree:"app_type"`
	InstanceType      value.String `tree:"instance_type"`
}
