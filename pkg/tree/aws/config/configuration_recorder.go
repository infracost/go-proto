package config

import "github.com/infracost/go-proto/pkg/tree/resource"

type ConfigurationRecorder struct {
	resource.Resource `tree:"-"`
}
