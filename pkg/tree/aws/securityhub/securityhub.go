package securityhub

import "github.com/infracost/go-proto/pkg/tree/resource"

type SecurityHub struct {
	Accounts []Account `tree:"accounts"`
}

type Account struct {
	resource.Resource `tree:"-"`
}
