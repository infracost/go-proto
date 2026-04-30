package secretsmanager

import "github.com/infracost/go-proto/pkg/tree/resource"

type Secret struct {
	resource.Resource `tree:"-"`
}
