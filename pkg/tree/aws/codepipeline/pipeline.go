package codepipeline

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Pipeline struct {
	resource.Resource `tree:"-"`
	PipelineType      value.Value[PipelineType] `tree:"pipeline_type"`
}

type PipelineType uint32

const (
	PipelineTypeUnknown PipelineType = iota
	PipelineTypeV1
	PipelineTypeV2
)
