package backup

import "github.com/infracost/go-proto/pkg/tree/resource"

type LogicallyAirGappedVault struct {
	resource.Resource `tree:"-"`
}
