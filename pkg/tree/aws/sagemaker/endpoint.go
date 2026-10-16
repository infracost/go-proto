package sagemaker

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Endpoint struct {
	resource.Resource  `tree:"-"`
	EndpointConfigName value.String `tree:"endpoint_config_name"`

	Relationships EndpointRelationships `tree:"-"`
}

type EndpointRelationships struct {
	EndpointConfiguration *EndpointConfiguration
}
