package config

import "github.com/infracost/go-proto/pkg/tree/resource"

type ConfigRule struct {
	resource.Resource `tree:"-"`
}
