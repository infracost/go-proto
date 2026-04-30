package automation

import "github.com/infracost/go-proto/pkg/tree/resource"

type Account struct {
	resource.Resource `tree:"-"`
}
