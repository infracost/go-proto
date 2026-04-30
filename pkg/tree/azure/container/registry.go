package container

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Registry struct {
	resource.Resource         `tree:"-"`
	SKU                       value.Value[RegistrySKU] `tree:"sku"`
	GeoReplicationLocations   value.List[string] `tree:"geo_replication_locations"`
	RetentionPolicyDays       value.Int          `tree:"retention_policy_days"`
}
