package monitor

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

type DataCollectionRule struct {
	resource.Resource `tree:"-"`
}
