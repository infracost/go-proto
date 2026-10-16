package datasync

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Task struct {
	resource.Resource `tree:"-"`
	TaskMode          value.Value[TaskMode] `tree:"task_mode"`
}

type TaskMode uint32

const (
	TaskModeUnknown TaskMode = iota
	TaskModeBasic
	TaskModeEnhanced
)
