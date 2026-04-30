package glue

import "github.com/infracost/go-proto/pkg/tree/resource"

type Crawler struct {
	resource.Resource `tree:"-"`
}
