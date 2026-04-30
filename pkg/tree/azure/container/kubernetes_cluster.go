package container

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type KubernetesCluster struct {
	resource.Resource   `tree:"-"`
	SKUTier             value.Value[KubernetesSKUTier] `tree:"sku_tier"`
	DefaultNodePool     NodePool     `tree:"default_node_pool"`
}
