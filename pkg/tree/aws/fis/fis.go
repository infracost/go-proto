package fis

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FIS struct {
	ExperimentTemplates []ExperimentTemplate `tree:"experiment_templates"`
}

type ExperimentTemplate struct {
	resource.Resource `tree:"-"`
	MultiAccount      value.Bool `tree:"multi_account"`
	HasReport         value.Bool `tree:"has_report"`
}
