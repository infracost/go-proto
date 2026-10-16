package sagemaker

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EndpointConfiguration struct {
	resource.Resource  `tree:"-"`
	Name               value.String      `tree:"name"`
	ProductionVariants []EndpointVariant `tree:"production_variants"`
	ShadowVariants     []EndpointVariant `tree:"shadow_production_variants"`
}

type EndpointVariant struct {
	VariantName          value.String `tree:"variant_name"`
	InstanceType         value.String `tree:"instance_type"`
	InitialInstanceCount value.Int    `tree:"initial_instance_count"`
	VolumeSizeInGB       value.Int    `tree:"volume_size_in_gb"`
	// Serverless variants have no instance type.
	ServerlessMemorySizeInMB         value.Int `tree:"serverless_memory_size_in_mb"`
	ServerlessProvisionedConcurrency value.Int `tree:"serverless_provisioned_concurrency"`
}
