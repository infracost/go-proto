package elasticache

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ServerlessCache struct {
	resource.Resource `tree:"-"`
	Engine            value.Value[CacheEngine] `tree:"engine"`
}
