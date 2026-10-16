package detective

import "github.com/infracost/go-proto/pkg/tree/resource"

type Detective struct {
	Graphs []Graph `tree:"graphs"`
}

type Graph struct {
	resource.Resource `tree:"-"`
}
