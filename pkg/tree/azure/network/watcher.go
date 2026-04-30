package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

type Watcher struct {
	resource.Resource `tree:"-"`
}
