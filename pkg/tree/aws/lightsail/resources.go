package lightsail

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Database struct {
	resource.Resource `tree:"-"`
	BundleID          value.String `tree:"bundle_id"`
}

type Disk struct {
	resource.Resource `tree:"-"`
	SizeGB            value.Int `tree:"size_gb"`
}

type LoadBalancer struct {
	resource.Resource `tree:"-"`
}

type ContainerService struct {
	resource.Resource `tree:"-"`
	Power             value.String `tree:"power"`
	Scale             value.Int    `tree:"scale"`
}

type Distribution struct {
	resource.Resource `tree:"-"`
	BundleID          value.String `tree:"bundle_id"`
}

type Bucket struct {
	resource.Resource `tree:"-"`
	BundleID          value.String `tree:"bundle_id"`
}
