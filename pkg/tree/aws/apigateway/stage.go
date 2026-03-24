package apigateway

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Stage struct {
	resource.Resource `tree:"-"`
	CacheClusterSize  value.String `tree:"cache_cluster_size"`
	CacheEnabled      value.Bool   `tree:"cache_enabled"`
}
